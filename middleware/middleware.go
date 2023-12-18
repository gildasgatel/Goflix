package middleware

import (
	"fmt"
	"goflix/models"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = []byte("secret_key")

func GenerateToken(id int, account string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["user_account"] = account
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token valide pendant 24 heures

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}
		user, err := ExtractUserData(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Data token invalide"})
			c.Abort()
			return
		}
		if user.Account != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "acces denied, admin only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ExtractUserData(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token invalide")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("impossible de récupérer les claims")
	}
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id n'est pas présent ou n'est pas de type int64")
	}
	userAccount, ok := claims["user_account"].(string)
	if !ok {
		return nil, fmt.Errorf("user_account n'est pas présent ou n'est pas de type string")
	}

	return &models.User{Id: int(userID), Account: userAccount}, nil
}
