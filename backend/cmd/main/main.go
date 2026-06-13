package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cinema-booking/config"
	"cinema-booking/internal/admin"
	"cinema-booking/internal/auth"
	"cinema-booking/internal/booking"
	"cinema-booking/internal/middleware"
	"cinema-booking/internal/movie"
	"cinema-booking/internal/queue"
	"cinema-booking/internal/websocket"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	db := config.ConnectMongo(cfg)
	rdb := config.ConnectRedis(cfg)
	var mq *queue.Publisher
    var err error 
	for i := 0; i < 10; i++ {
		mq, err = queue.NewPublisher(cfg.RabbitMQURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ not ready, retrying in 5s... (%d/10)", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("RabbitMQ connection failed after retries: %v", err)
	}

	go queue.StartConsumer(cfg.RabbitMQURL, db)

	hub := websocket.NewHub()
	go hub.Run()
    go booking.StartTimeoutWatcher(db, rdb, hub)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(middleware.CORS())

	r.POST("/api/auth/google", auth.GoogleLogin(db, cfg))
	r.GET("/api/movies", movie.List(db))

	r.GET("/ws", websocket.Handler(hub))

	api := r.Group("/api", middleware.AuthRequired(cfg))
	{
		api.GET("/movies/:id/showtimes", movie.Showtimes(db))
		api.GET("/showtimes/:id/seats", booking.GetSeats(db, rdb))
		api.POST("/bookings", booking.Create(db, rdb, hub, mq))
		api.POST("/bookings/:id/payment", booking.ConfirmPayment(db, rdb, hub, mq))
		api.GET("/bookings/my", booking.MyBookings(db))
	}

	adminGroup := r.Group("/api/admin",
		middleware.AuthRequired(cfg),
		middleware.AdminOnly(),
	)
	{
		adminGroup.GET("/bookings", admin.ListBookings(db))
		adminGroup.GET("/audit-logs", admin.AuditLogs(db))
        adminGroup.GET("/stats",      admin.Stats(db))
	}

	r.GET("/health", func(c *gin.Context) {
		if err := db.Client().Ping(context.Background(), nil); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "mongo": "down"})
			return
		}
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Server running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
