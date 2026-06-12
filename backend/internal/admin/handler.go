package admin

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ListBookings(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        movieID  := c.Query("movie_id")   
        status   := c.Query("status")     
        dateStr  := c.Query("date")       

        filter := bson.M{}

        if status != "" {
            filter["status"] = status
        }

        if dateStr != "" {
            date, err := time.Parse("2006-01-02", dateStr)
            if err == nil {
                filter["created_at"] = bson.M{
                    "$gte": date,
                    "$lt":  date.Add(24 * time.Hour),
                }
            }
        }

        var pipeline mongo.Pipeline

        if movieID != "" {
            oid, _ := primitive.ObjectIDFromHex(movieID)
            pipeline = mongo.Pipeline{
                {{Key: "$lookup", Value: bson.M{
                    "from":         "showtimes",
                    "localField":   "showtime_id",
                    "foreignField": "_id",
                    "as":           "showtime",
                }}},
                {{Key: "$match", Value: bson.M{
                    "showtime.movie_id": oid,
                }}},
                {{Key: "$lookup", Value: bson.M{
                    "from":         "users",
                    "localField":   "user_id",
                    "foreignField": "_id",
                    "as":           "user",
                }}},
                {{Key: "$sort", Value: bson.M{"created_at": -1}}},
                {{Key: "$limit", Value: 100}},
            }
        } else {
            pipeline = mongo.Pipeline{
                {{Key: "$match", Value: filter}},
                {{Key: "$lookup", Value: bson.M{
                    "from":         "users",
                    "localField":   "user_id",
                    "foreignField": "_id",
                    "as":           "user",
                }}},
                {{Key: "$lookup", Value: bson.M{
                    "from":         "showtimes",
                    "localField":   "showtime_id",
                    "foreignField": "_id",
                    "as":           "showtime",
                }}},
                {{Key: "$sort",  Value: bson.M{"created_at": -1}}},
                {{Key: "$limit", Value: 100}},
            }
        }

        cursor, err := db.Collection("bookings").Aggregate(
            context.Background(), pipeline,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer cursor.Close(context.Background())

        var results []bson.M
        if err := cursor.All(context.Background(), &results); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "bookings": results,
            "total":    len(results),
        })
    }
}

func AuditLogs(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        eventType := c.Query("event")  
        userID    := c.Query("user_id")

        filter := bson.M{}
        if eventType != "" {
            filter["event"] = eventType
        }
        if userID != "" {
            filter["user_id"] = userID
        }

        opts := options.Find().
            SetSort(bson.M{"created_at": -1}). 
            SetLimit(200)

        cursor, err := db.Collection("audit_logs").Find(
            context.Background(), filter, opts,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        defer cursor.Close(context.Background())

        var logs []bson.M
        cursor.All(context.Background(), &logs)

        c.JSON(http.StatusOK, gin.H{"logs": logs})
    }
}

func Stats(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := context.Background()

        pipeline := mongo.Pipeline{
            {{Key: "$group", Value: bson.M{
                "_id":   "$status",
                "count": bson.M{"$sum": 1},
            }}},
        }

        cursor, _ := db.Collection("bookings").Aggregate(ctx, pipeline)
        var stats []bson.M
        cursor.All(ctx, &stats)

        userCount, _ := db.Collection("users").CountDocuments(ctx, bson.M{})

        c.JSON(http.StatusOK, gin.H{
            "booking_stats": stats,
            "total_users":   userCount,
        })
    }
}