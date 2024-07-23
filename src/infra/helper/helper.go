package helper

import (
	dto "e-depo/src/app/dto/user"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret_key")

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrUserNotFound        = errors.New("destination user not found")
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
