package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	CreateClaims(userID int64) jwt.Claims
	CreateToken(claims jwt.Claims) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
	exp    time.Duration
}

func NewJWTAuthenticator(secret, aud, iss string, exp time.Duration) *JWTAuthenticator {
	return &JWTAuthenticator{secret, iss, aud, exp}
}

func (a *JWTAuthenticator) CreateClaims(userID int64) jwt.Claims {
	return jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.exp)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    a.iss,
		Subject:   strconv.FormatInt(userID, 10),
		Audience:  []string{a.aud},
	}
}

func (a *JWTAuthenticator) CreateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}

func (a *JWTAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
