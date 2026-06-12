package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type AuditLog struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
    Event     string                 `bson:"event" json:"event"`
    UserID    string                 `bson:"user_id" json:"user_id"`
    Data      map[string]interface{} `bson:"data" json:"data"`  
    CreatedAt time.Time              `bson:"created_at" json:"created_at"`
}