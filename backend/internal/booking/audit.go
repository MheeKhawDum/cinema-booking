package booking

import (
    "context"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "cinema-booking/internal/models"
)

func logAudit(db *mongo.Database, event, userID string, data map[string]interface{}) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    db.Collection("audit_logs").InsertOne(ctx, models.AuditLog{
        ID:        primitive.NewObjectID(),
        Event:     event,
        UserID:    userID,
        Data:      data,
        CreatedAt: time.Now(),
    })
}