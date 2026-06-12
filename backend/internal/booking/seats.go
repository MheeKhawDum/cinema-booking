package booking

import (
    "context"
    "net/http"
    "time"
	"fmt"

    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

type SeatResponse struct {
    Number string `json:"number"`
    Row    string `json:"row"`    
    Col    int    `json:"col"`    
    Status string `json:"status"` 
}

func GetSeats(db *mongo.Database, rdb *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        showtimeID, err := primitive.ObjectIDFromHex(c.Param("id"))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid showtime id"})
            return
        }

        cursor, err := db.Collection("bookings").Find(ctx, bson.M{
            "showtime_id": showtimeID,
            "status":      "BOOKED",
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
            return
        }
        defer cursor.Close(ctx)

        bookedSeats := map[string]bool{}
        for cursor.Next(ctx) {
            var b struct {
                SeatNumber string `bson:"seat_number"`
            }
            if cursor.Decode(&b) == nil {
                bookedSeats[b.SeatNumber] = true
            }
        }

        rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
        var seats []SeatResponse

        for _, row := range rows {
            for col := 1; col <= 10; col++ {
                seatNum := fmt.Sprintf("%s%d", row, col)
                status := "AVAILABLE"

                if bookedSeats[seatNum] {
                    status = "BOOKED"
                } else {
                    key := fmt.Sprintf("lock:seat:%s:%s", showtimeID.Hex(), seatNum)
                    exists, err := rdb.Exists(ctx, key).Result()
                    if err == nil && exists > 0 {
                        status = "LOCKED"
                    }
                }

                seats = append(seats, SeatResponse{
                    Number: seatNum,
                    Row:    row,
                    Col:    col,
                    Status: status,
                })
            }
        }

        c.JSON(http.StatusOK, gin.H{
            "showtime_id": showtimeID.Hex(),
            "seats":       seats,
            "summary": gin.H{
                "total":     len(seats),
                "available": countByStatus(seats, "AVAILABLE"),
                "locked":    countByStatus(seats, "LOCKED"),
                "booked":    countByStatus(seats, "BOOKED"),
            },
        })
    }
}

func countByStatus(seats []SeatResponse, status string) int {
    n := 0
    for _, s := range seats {
        if s.Status == status {
            n++
        }
    }
    return n
}