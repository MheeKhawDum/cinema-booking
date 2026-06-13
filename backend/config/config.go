package config

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/redis/go-redis/v9"
)

type Config struct {
	MongoURI        string
    RedisAddr       string
    RedisPassword   string
    RabbitMQURL     string
    JWTSecret       string
    GoogleClientID  string
    GoogleSecret    string
}

func Load() *Config {
    return &Config{
        MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
        RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
        RedisPassword:  getEnv("REDIS_PASSWORD", ""),
        RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
        JWTSecret:      getEnv("JWT_SECRET", "dev-secret"),
        GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
        GoogleSecret:   getEnv("GOOGLE_CLIENT_SECRET", ""),
    }
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}

func ConnectMongo(cfg *Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("🔌 Connecting to MongoDB: %s", cfg.MongoURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("MongoDB connect error: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping error: %v", err)
	}

	log.Println("✅ MongoDB connected")
	return client.Database("cinema")
}
func ConnectRedis(cfg *Config) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.RedisAddr,
        Password: cfg.RedisPassword,
        DB:       0, 
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if _, err := rdb.Ping(ctx).Result(); err != nil {
        log.Fatalf("Redis connect error: %v", err)
    }

    log.Println("✅ Redis connected")
    return rdb
}
