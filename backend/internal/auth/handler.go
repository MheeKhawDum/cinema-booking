package auth

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "google.golang.org/api/idtoken"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

	"cinema-booking/config"   
    "cinema-booking/internal/models" 
)

type GoogleLoginRequest struct {
    IDToken string `json:"id_token" binding:"required"`
}

func GoogleLogin(db *mongo.Database, cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req GoogleLoginRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "id_token required"})
            return
        }

        payload, err := idtoken.Validate(context.Background(), req.IDToken, cfg.GoogleClientID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid google token"})
            return
        }

        googleID := payload.Subject                    
        email    := payload.Claims["email"].(string)
        name     := payload.Claims["name"].(string)
        picture  := payload.Claims["picture"].(string)

        coll := db.Collection("users")
        filter := bson.M{"google_id": googleID}
        update := bson.M{
            "$set": bson.M{
                "email":      email,
                "name":       name,
                "picture":    picture,
                "updated_at": time.Now(),
            },
            "$setOnInsert": bson.M{
                "google_id":  googleID,
                "role":       "USER",          
                "created_at": time.Now(),
            },
        }
        opts := options.FindOneAndUpdate().
            SetUpsert(true).
            SetReturnDocument(options.After)   

        var user models.User
        err = coll.FindOneAndUpdate(context.Background(), filter, update, opts).
            Decode(&user)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
            return
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
            "user_id": user.ID.Hex(),
            "role":    user.Role,
            "exp":     time.Now().Add(24 * time.Hour).Unix(),
        })

        tokenStr, err := token.SignedString([]byte(cfg.JWTSecret))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "token": tokenStr,
            "user":  user,
        })
    }
}