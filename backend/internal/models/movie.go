package models

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string             `bson:"title"         json:"title"`
    Description string             `bson:"description"   json:"description"`
    PosterURL   string             `bson:"poster_url"    json:"poster_url"`
    Duration    int                `bson:"duration"      json:"duration"`
    CreatedAt   time.Time          `bson:"created_at"    json:"created_at"`
}

type Showtime struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    MovieID    primitive.ObjectID `bson:"movie_id"      json:"movie_id"`
    StartTime  time.Time          `bson:"start_time"    json:"start_time"`
    Hall       string             `bson:"hall"          json:"hall"`
    TotalSeats int                `bson:"total_seats"   json:"total_seats"`
}