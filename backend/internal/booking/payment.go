package booking

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/redis/go-redis/v9"

    "cinema-booking/internal/models"
    "cinema-booking/internal/queue"
    "cinema-booking/internal/websocket"
)

type PaymentRequest struct {
    BookingID string `json:"booking_id" binding:"required"`
}

func ConfirmPayment(db *mongo.Database, rdb *redis.Client, hub *websocket.Hub, mq *queue.Publisher) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")

        var req PaymentRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        ctx := context.Background()

        var b models.Booking
        err := db.Collection("bookings").FindOne(ctx, bson.M{
            "_id":     objectIDFromString(req.BookingID),
            "user_id": objectIDFromString(userID.(string)),
            "status":  "PENDING",
        }).Decode(&b)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "booking not found or expired"})
            return
        }

        if time.Now().After(b.ExpiredAt) {
            ReleaseLock(ctx, rdb, b.ShowtimeID.Hex(), b.SeatNumber, userID.(string))
            db.Collection("bookings").UpdateOne(ctx,
                bson.M{"_id": b.ID},
                bson.M{"$set": bson.M{"status": "TIMEOUT"}},
            )
            hub.Broadcast(websocket.SeatEvent{
                ShowtimeID: b.ShowtimeID.Hex(),
                Seat:       b.SeatNumber,
                Status:     "AVAILABLE",
            })
            c.JSON(http.StatusGone, gin.H{"error": "payment timeout, seat released"})
            return
        }

        now := time.Now()
        db.Collection("bookings").UpdateOne(ctx,
            bson.M{"_id": b.ID},
            bson.M{"$set": bson.M{"status": "BOOKED", "booked_at": now}},
        )

        ReleaseLock(ctx, rdb, b.ShowtimeID.Hex(), b.SeatNumber, userID.(string))

        hub.Broadcast(websocket.SeatEvent{
            ShowtimeID: b.ShowtimeID.Hex(),
            Seat:       b.SeatNumber,
            Status:     "BOOKED",
        })

        mq.Publish("booking.success", map[string]interface{}{
            "booking_id": req.BookingID,
            "user_id":    userID.(string),
            "seat":       b.SeatNumber,
        })

        logAudit(db, "BOOKING_SUCCESS", userID.(string), map[string]interface{}{
            "booking_id": req.BookingID,
            "seat":       b.SeatNumber,
        })

        c.JSON(http.StatusOK, gin.H{
            "message":    "booking confirmed!",
            "booking_id": req.BookingID,
        })
    }
}