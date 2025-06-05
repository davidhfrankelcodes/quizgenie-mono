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
// (Not strictly required if you use MapClaims, but provided here for illustration.)
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

	// Set standard + custom claims
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken validates the given JWT string and extracts the user ID.
// It returns ErrInvalidToken if parsing/validation fails, or ErrNotInitialized if InitJWT was not called.
func ParseToken(tokenString string) (uint, error) {
	if jwtSecret == nil {
		return 0, ErrNotInitialized
	}

	parsed, err := jwt.Parse(tokenString, func(tok *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return 0, ErrInvalidToken
	}

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		// The "sub" claim is stored as float64 by default; convert to uint
		subRaw, exists := claims["sub"]
		if !exists {
			return 0, ErrInvalidToken
		}
		idFloat, ok := subRaw.(float64)
		if !ok {
			return 0, ErrInvalidToken
		}
		return uint(idFloat), nil
	}

	return 0, ErrInvalidToken
}
