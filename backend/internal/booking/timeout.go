package booking

import (
    "context"
    "log"
    "time"
	"fmt"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/redis/go-redis/v9"

    "cinema-booking/internal/models"
    "cinema-booking/internal/websocket"
)

func StartTimeoutWatcher(db *mongo.Database, rdb *redis.Client, hub *websocket.Hub) {
    log.Println("⏱  Timeout watcher started")

    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        checkExpiredBookings(db, rdb, hub)
    }
}

func checkExpiredBookings(db *mongo.Database, rdb *redis.Client, hub *websocket.Hub) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    cursor, err := db.Collection("bookings").Find(ctx, bson.M{
        "status":     "PENDING",
        "expired_at": bson.M{"$lt": time.Now()},
    })
    if err != nil {
        log.Printf("Timeout watcher error: %v", err)
        return
    }
    defer cursor.Close(ctx)

    var expired []models.Booking
    cursor.All(ctx, &expired)

    if len(expired) == 0 {
        return
    }

    log.Printf("⏱  Found %d expired bookings", len(expired))

    for _, booking := range expired {
        _, err := db.Collection("bookings").UpdateOne(ctx,
            bson.M{"_id": booking.ID},
            bson.M{"$set": bson.M{
                "status":     "TIMEOUT",
                "updated_at": time.Now(),
            }},
        )
        if err != nil {
            log.Printf("Error updating booking %s: %v", booking.ID.Hex(), err)
            continue
        }

        lockKey := fmt.Sprintf("lock:seat:%s:%s",
            booking.ShowtimeID.Hex(), booking.SeatNumber)
        rdb.Del(ctx, lockKey)

        hub.Broadcast(websocket.SeatEvent{
            ShowtimeID: booking.ShowtimeID.Hex(),
            Seat:       booking.SeatNumber,
            Status:     "AVAILABLE",
        })

        db.Collection("audit_logs").InsertOne(ctx, models.AuditLog{
            ID:    primitive.NewObjectID(),
            Event: "BOOKING_TIMEOUT",
            UserID: booking.UserID.Hex(),
            Data: map[string]interface{}{
                "booking_id": booking.ID.Hex(),
                "seat":       booking.SeatNumber,
                "showtime_id": booking.ShowtimeID.Hex(),
            },
            CreatedAt: time.Now(),
        })

        log.Printf("✅ Booking %s timed out, seat %s released",
            booking.ID.Hex(), booking.SeatNumber)
    }
}