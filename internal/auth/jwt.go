package auth

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret          string
	AccessIssuer    string
	RefreshIssuer   string
	AccessDuration  time.Duration
	RefreshDuration time.Duration
}

func (cfg JWTConfig) CreateToken(id int, expiresIn time.Duration, issuer string) (string, error) {
	now := time.Now().UTC()
	claims := &jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   strconv.Itoa(id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.Secret))
	return ss, err
}

func (cfg JWTConfig) ValidateToken(authToken string) (int, string, error) {
	ss := strings.TrimPrefix(authToken, "Bearer ")
	token, err := jwt.ParseWithClaims(ss, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return 0, "", err
	}
	stringId, err := token.Claims.GetSubject()
	if err != nil {
		return 0, "", err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return 0, "", err
	}
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return 0, "", err
	}
	return id, issuer, nil
}
