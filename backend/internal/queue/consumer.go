package queue

import (
    "context"
    "encoding/json"
    "log"
    "time"
	"fmt"  

    amqp "github.com/rabbitmq/amqp091-go"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "cinema-booking/internal/models"
)

func StartConsumer(url string, db *mongo.Database) {
    conn, err := amqp.Dial(url)
    if err != nil {
        log.Printf("Consumer connect error: %v", err)
        return
    }
    ch, _ := conn.Channel()
    ch.QueueDeclare("booking.events", true, false, false, false, nil)
    msgs, _ := ch.Consume("booking.events", "", true, false, false, false, nil)

    log.Println("📨 Consumer started")

    for msg := range msgs {
        var event map[string]interface{}
        json.Unmarshal(msg.Body, &event)

        db.Collection("audit_logs").InsertOne(
            context.Background(),
            models.AuditLog{
                ID:        primitive.NewObjectID(),
                Event:     "BOOKING_SUCCESS",
                UserID:    fmt.Sprintf("%v", event["user_id"]),
                Data:      event,
                CreatedAt: time.Now(),
            },
        )
        log.Printf("📧 Mock notification sent: %v", event)
    }
}