package middleware

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	UserId uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func generateSecretKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func CreateToken(userId uint, role string, secretKey string) (string, error) {
	claims := CustomClaims{
		userId,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "myapp",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func AuthMiddleware(secretKey string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == " " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString[7:], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !contains(allowedRoles, claims.Role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.Set("userId", claims.UserId)
		c.Next()
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func SetCookie(c *gin.Context, name string, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(c.Writer, cookie)
}

func GetUserIdAndRoleFromCookie(c *gin.Context, secretKey string) (uint, string, error) {
	// Get JWT token from cookie
	cookie, err := c.Cookie("jwt-token")
	if err != nil {
		return 0, "", err
	}

	// Extract token from "Bearer <token>" format
	tokenString := strings.Replace(cookie, "Bearer ", "", 1)

	// Parse JWT token and extract user ID and role
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}

		// Return secret key as signing key
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", jwt.NewValidationError("Invalid JWT token", jwt.ValidationErrorSignatureInvalid)
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", jwt.NewValidationError("Invalid user ID in JWT token", jwt.ValidationErrorMalformed)
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", jwt.NewValidationError("Invalid role in JWT token", jwt.ValidationErrorMalformed)
	}

	return uint(userId), role, nil
}
