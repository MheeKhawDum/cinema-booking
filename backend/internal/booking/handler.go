package booking

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/redis/go-redis/v9"

    "cinema-booking/internal/models"
    "cinema-booking/internal/queue"
    "cinema-booking/internal/websocket"
)

type BookingRequest struct {
    ShowtimeID string `json:"showtime_id" binding:"required"`
    Seat       string `json:"seat"        binding:"required"`
}

func objectIDFromString(id string) primitive.ObjectID {
    oid, _ := primitive.ObjectIDFromHex(id)
    return oid
}

func checkBookingExists(ctx context.Context, db *mongo.Database, showtimeID, seat string) (bool, error) {
    count, err := db.Collection("bookings").CountDocuments(ctx, bson.M{
        "showtime_id": objectIDFromString(showtimeID),
        "seat_number": seat,
        "status":      "BOOKED",
    })
    return count > 0, err
}

func Create(db *mongo.Database, rdb *redis.Client, hub *websocket.Hub, mq *queue.Publisher) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")

        var req BookingRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        ctx := context.Background()

        exists, err := checkBookingExists(ctx, db, req.ShowtimeID, req.Seat)
        if err != nil || exists {
            c.JSON(http.StatusConflict, gin.H{"error": "seat already booked"})
            return
        }

        acquired, err := AcquireLock(ctx, rdb, req.ShowtimeID, req.Seat, userID.(string))
        if err != nil {
            logAudit(db, "LOCK_FAIL", userID.(string), map[string]interface{}{
                "seat": req.Seat, "error": err.Error(),
            })
            c.JSON(http.StatusInternalServerError, gin.H{"error": "lock error"})
            return
        }

        if !acquired {
            c.JSON(http.StatusConflict, gin.H{"error": "seat is being selected by someone"})
            return
        }

        hub.Broadcast(websocket.SeatEvent{
            ShowtimeID: req.ShowtimeID,
            Seat:       req.Seat,
            Status:     "LOCKED",
            UserID:     userID.(string),
        })

        b := models.Booking{
            UserID:     objectIDFromString(userID.(string)),
            ShowtimeID: objectIDFromString(req.ShowtimeID),
            SeatNumber: req.Seat,
            Status:     "PENDING",
            Amount:     150,
            ExpiredAt:  time.Now().Add(5 * time.Minute),
            CreatedAt:  time.Now(),
        }
        result, err := db.Collection("bookings").InsertOne(ctx, b)
        if err != nil {
            ReleaseLock(ctx, rdb, req.ShowtimeID, req.Seat, userID.(string))
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
            return
        }

        bookingID := result.InsertedID.(primitive.ObjectID).Hex()

        c.JSON(http.StatusOK, gin.H{
            "booking_id": bookingID,
            "expires_at": b.ExpiredAt,
            "message":    "seat locked, please complete payment within 5 minutes",
        })
    }
}

func MyBookings(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")
        ctx := context.Background()

        oid := objectIDFromString(userID.(string))

        pipeline := mongo.Pipeline{
            {{Key: "$match", Value: bson.M{"user_id": oid}}},
            {{Key: "$lookup", Value: bson.M{
                "from": "showtimes", "localField": "showtime_id",
                "foreignField": "_id", "as": "showtime",
            }}},
            {{Key: "$unwind", Value: bson.M{"path": "$showtime", "preserveNullAndEmptyArrays": true}}},
            {{Key: "$lookup", Value: bson.M{
                "from": "movies", "localField": "showtime.movie_id",
                "foreignField": "_id", "as": "movie",
            }}},
            {{Key: "$unwind", Value: bson.M{"path": "$movie", "preserveNullAndEmptyArrays": true}}},
            {{Key: "$sort", Value: bson.M{"created_at": -1}}},
        }

        cursor, err := db.Collection("bookings").Aggregate(ctx, pipeline)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
            return
        }
        defer cursor.Close(ctx)

        var bookings []bson.M
        cursor.All(ctx, &bookings)
        if bookings == nil {
            bookings = []bson.M{}
        }
        c.JSON(http.StatusOK, gin.H{"bookings": bookings})
    }
}