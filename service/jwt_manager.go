package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWTManager manages jwt token services
type JWManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// UserClaims contains information about the user
type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

// NewJWTManager return a new jwt manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWManager {
	return &JWManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// Generate generates and returns a jwt token for the user
func (m *JWManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// Verify verifies the token and returns user's claims
func (m *JWManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(m.secretKey), nil
		})

	if err != nil {
		return nil, fmt.Errorf("incalid token: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
