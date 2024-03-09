package jwt

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"gorm.io/gorm"
	"job_board/db"
	"job_board/helpers"
	"job_board/models"
)

var database *gorm.DB
var SecretKey []byte

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

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			helpers.CreateResponse(c, helpers.Response{
				Message:    "invalid jwt",
				StatusCode: http.StatusUnauthorized,
				Data:       nil,
			})
			// c.String(http.StatusUnauthorized, "invalid jwt")
			// c.Abort()
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
			// c.String(http.StatusBadRequest, err.Error())
			// c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			helpers.CreateResponse(c, helpers.Response{
				Message:    "invalid jwt",
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			// c.String(http.StatusBadRequest, "invalid jwt")
			// c.Abort()
			return
		}

		providerID := claims["provider_id"].(string)
		session := sessions.Default(c)
		val, ok := session.Get(providerID).(models.User)
		if ok {
			c.Set("claims", claims)
			c.Set("user", val)
			c.Next()
			return
		}

		var user models.User
		if err := database.Preload("Profile").Preload("Companies").Preload("JobApplications").Where(&models.User{ProviderID: providerID}).First(&user).Error; err != nil {
			helpers.CreateResponse(c, helpers.Response{
				Message:    err.Error(),
				StatusCode: http.StatusBadRequest,
				Data:       nil,
			})
			// c.String(http.StatusBadRequest, err.Error())
			// c.Abort()
			return
		}
		session.Set(providerID, user)
		c.Set("claims", claims)
		c.Set("user", user)
		c.Next()
	}
}
