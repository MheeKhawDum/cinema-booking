package movie

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "cinema-booking/internal/models"
)

func List(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        opts := options.Find().SetSort(bson.M{"title": 1})
        cursor, err := db.Collection("movies").Find(ctx, bson.M{}, opts)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
            return
        }
        defer cursor.Close(ctx)

        var movies []models.Movie
        if err := cursor.All(ctx, &movies); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "decode error"})
            return
        }

        if movies == nil {
            movies = []models.Movie{}
        }

        c.JSON(http.StatusOK, gin.H{"movies": movies})
    }
}

func Showtimes(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        movieID, err := primitive.ObjectIDFromHex(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
            return
        }

        var movie models.Movie
        err = db.Collection("movies").
            FindOne(ctx, bson.M{"_id": movieID}).
            Decode(&movie)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
            return
        }

        opts := options.Find().SetSort(bson.M{"start_time": 1})
        cursor, err := db.Collection("showtimes").Find(ctx,
            bson.M{"movie_id": movieID},
            opts,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
            return
        }
        defer cursor.Close(ctx)

        var showtimes []models.Showtime
        if err := cursor.All(ctx, &showtimes); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "decode error"})
            return
        }

        if showtimes == nil {
            showtimes = []models.Showtime{}
        }

        c.JSON(http.StatusOK, gin.H{
            "movie":     movie,
            "showtimes": showtimes,
        })
    }
}