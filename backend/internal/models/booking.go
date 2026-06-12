package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	
	"time"
)

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
    ShowtimeID  primitive.ObjectID `bson:"showtime_id" json:"showtime_id"`
    SeatNumber  string             `bson:"seat_number" json:"seat_number"` 
    Status      string             `bson:"status" json:"status"`
    Amount      float64            `bson:"amount" json:"amount"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    ExpiredAt   time.Time          `bson:"expired_at" json:"expired_at"` 
}