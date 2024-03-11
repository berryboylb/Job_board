package jwt

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
	"job_board/db"
	"job_board/helpers"
	"job_board/models"
	cisredis "job_board/redis"
)

var SecretKey []byte
var database *gorm.DB

func init() {
	database = db.GetDB()
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("Error loading jwt secret")
	}
	SecretKey = []byte(secretKey)
}

func GenerateJWT(providerID string, isMobile bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["provider_id"] = providerID

	if isMobile {
		// Set expiration time to 1 year from now for mobile users
		claims["exp"] = time.Now().Add(time.Hour * 24 * 365).Unix()
	} else {
		// Set expiration time to 30 minutes from now for web users
		claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	}

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetSingleUser(filter models.User) (*models.User, error) {
	var user models.User
	if err := database.
		Preload("Profile").
		Preload("Profile").
		Preload("Companies").
		Preload("JobApplications").
		Where(&filter).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUser(providerID string) (models.User, error) {
	userStr, err := cisredis.Retrieve(providerID)
	if err != nil {
		if err == redis.Nil {
			// Fetch user data from database
			user, err := GetSingleUser(models.User{ProviderID: providerID})
			if err != nil {
				return models.User{}, fmt.Errorf("failed to fetch user data from database: %w", err)
			}
			// Store user data in Redis with an expiration time
			userStr, err = cisredis.StoreStruct(user) // Fixed: Remove redeclaration
			if err != nil {
				return models.User{}, fmt.Errorf("failed to store user data in Redis: %w", err)
			}
			expiration := 10 * time.Minute
			err = cisredis.Store(providerID, userStr, expiration)
			if err != nil {
				return models.User{}, fmt.Errorf("failed to store user data in Redis with expiration: %w", err)
			}
			return *user, nil
		}
		return models.User{}, fmt.Errorf("failed to retrieve user data from Redis: %w", err)
	}

	var user models.User
	if err := cisredis.UnmarshalStruct(userStr.([]byte), &user); err != nil {
		return models.User{}, fmt.Errorf("failed to unmarshal user data from Redis: %w", err)
	}
	return user, nil
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.CreateResponse(c, helpers.Response{
				Message:    "invalid jwt",
				StatusCode: http.StatusUnauthorized,
				Data:       nil,
			})
			return
		}

		token, err := jwt.Parse(authHeader[len("Bearer "):], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Return the secret key
			return SecretKey, nil
		})

		if err != nil {
			helpers.CreateResponse(c, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			helpers.CreateResponse(c, helpers.Response{
				Message:    "invalid jwt",
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			return
		}

		providerID := claims["provider_id"].(string)
		user, err := GetUser(providerID)
		if err != nil {
			helpers.CreateResponse(c, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusInternalServerError,
				Data:       nil,
			})
			return
		}
		c.Set("claims", claims)
		c.Set("user", user)
		c.Next()
	}
}
