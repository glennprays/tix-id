package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = os.Getenv("JWT_KEY")
var tokenName = "token"

type CustomClaims struct {
	UserId uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func CreateToken(c *gin.Context, userId uint, role string, expSecond int) {
	claims := CustomClaims{
		userId,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expSecond)).Unix(),
			Issuer:    "tix-id",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.Println(err)
	}

	setCookie(c, signedToken, time.Second*time.Duration(expSecond))
}

func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie(tokenName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Access token is missing, please login first"})
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
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

func setCookie(c *gin.Context, value string, exp time.Duration) {
	cookie := &http.Cookie{
		Name:     tokenName,
		Value:    value,
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(exp),
		Path:     "/",
	}
	http.SetCookie(c.Writer, cookie)
}
func ResetUserToken(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     tokenName,
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now(),
		Path:     "/",
	}
	http.SetCookie(c.Writer, cookie)
}

func GetUserIdAndRoleFromCookie(c *gin.Context) (uint, string, error) {
	// Get JWT token from cookie
	cookie, err := c.Cookie(tokenName)
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
		return []byte(jwtKey), nil
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
