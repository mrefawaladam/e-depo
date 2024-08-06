package helper

import (
	"context"
	dto "e-depo/src/app/dto/user"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret_key")

type ContextKey string

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrUserNotFound        = errors.New("destination user not found")
)

const (
	ContextUserKey ContextKey = "user"
)

// TokenClaims menyimpan klaim JWT
type TokenClaims struct {
	UserID      int64  `json:"user_id"`
	UserName    string `json:"username"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken membuat token JWT
func GenerateToken(data *dto.UserModel) (string, error) {
	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &TokenClaims{
		UserID:      data.ID,
		Role:        data.Role,
		UserName:    data.UserName,
		Name:        data.Name,
		PhoneNumber: data.PhoneNumber,
		Address:     data.Address,
		Email:       data.Email,

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// VerifyToken memverifikasi token JWT
func VerifyToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

// RoleCheckMiddleware verifies the user's role from the JWT token
func RoleCheckMiddleware(requiredRole string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := VerifyToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims.Role != requiredRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
