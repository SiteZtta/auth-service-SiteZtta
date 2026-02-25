package jwt

import (
	"auth-service-SiteZtta/config"
	"auth-service-SiteZtta/internal/domain/entities"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrParseClaims = errors.New("error parsing claims")
	ErrParseToken  = errors.New("error parsing token")
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId int64 `json:"user_id"`
	Role   int32 `json:"role"`
}

func NewToken(user entities.User, authConf config.AuthConf) (string, error) {
	expiresAt := jwt.NumericDate{Time: time.Now().Add(authConf.TokenTtl).UTC()}
	issuedAt := jwt.NumericDate{Time: time.Now().UTC()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &expiresAt,
			IssuedAt:  &issuedAt,
		},
		UserId: user.ID,
		Role:   user.Role,
	})
	tokenSigned, err := token.SignedString([]byte(authConf.SigningKey))
	if err != nil {
		return "", err
	}
	return tokenSigned, nil
}

func ParseToken(tokenString string, authConf config.AuthConf) (TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authConf.SigningKey), nil
	})
	if err != nil {
		return TokenClaims{}, ErrParseToken
	}
	claims, ok := parsedToken.Claims.(*TokenClaims)
	if !ok {
		return TokenClaims{}, ErrParseClaims
	}
	return *claims, nil
}
