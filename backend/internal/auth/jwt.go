// backend/internal/auth/jwt.go

package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret holds the HMAC secret key for signing JWTs.
// It must be initialized by calling InitJWT() before any GenerateToken/ParseToken calls.
var jwtSecret []byte

// ErrNotInitialized is returned when GenerateToken or ParseToken is called
// before InitJWT has been invoked (i.e., jwtSecret is still nil).
var ErrNotInitialized = errors.New("JWT module not initialized; call InitJWT first")

// ErrInvalidToken is returned when a token fails validation or cannot be parsed.
var ErrInvalidToken = errors.New("invalid or expired JWT token")

// InitJWT reads the JWT_SECRET environment variable and stores it for future use.
// If JWT_SECRET is not set, this function will log.Fatal.
func InitJWT() {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	jwtSecret = []byte(s)
}

// Claims defines the structure embedded in each JWT.
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT containing the user ID and username.
// It returns an error if InitJWT was not called first.
func GenerateToken(userID uint, username string) (string, error) {
	if jwtSecret == nil {
		return "", ErrNotInitialized
	}

	// (We no longer declare a local "claims" variable that goes unused.)

	// Use MapClaims so that "sub" and "username" appear exactly as expected
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtSecret)
}

// ParseToken validates the given JWT string and returns a *Claims struct.
// It returns ErrInvalidToken if parsing/validation fails, or ErrNotInitialized if InitJWT was not called.
func ParseToken(tokenString string) (*Claims, error) {
	if jwtSecret == nil {
		return nil, ErrNotInitialized
	}

	parsed, err := jwt.Parse(tokenString, func(tok *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	if mapClaims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		// Extract "sub" (user ID) from the map
		subRaw, exists := mapClaims["sub"]
		if !exists {
			return nil, ErrInvalidToken
		}
		idFloat, ok := subRaw.(float64)
		if !ok {
			return nil, ErrInvalidToken
		}
		userID := uint(idFloat)

		// Extract "username"
		usernameRaw, exists := mapClaims["username"]
		if !exists {
			return nil, ErrInvalidToken
		}
		username, ok := usernameRaw.(string)
		if !ok {
			return nil, ErrInvalidToken
		}

		// Build and return a Claims object
		return &Claims{
			UserID:   userID,
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(int64(mapClaims["exp"].(float64)), 0)),
			},
		}, nil
	}

	return nil, ErrInvalidToken
}
