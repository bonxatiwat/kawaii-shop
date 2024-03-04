package kawaiiauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/bonxatiwat/kawaii-shop-tutortial/config"
	"github.com/bonxatiwat/kawaii-shop-tutortial/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type kawaiiAuth struct {
	mapClaims *kawaiiMapCliams
	cfg       config.IJwtConfig
}

type kawaiiAdmin struct {
	*kawaiiAuth
}
type kawaiiMapCliams struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IKawaiiAuth interface {
	SignToken() string
}

type IKwaaiiAdmin interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate((time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9))))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *kawaiiAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *kawaiiAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func PassToken(cfg config.IJwtConfig, tokenString string) (*kawaiiMapCliams, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kawaiiMapCliams{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("singing method is invaild")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invaild")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expried")
		} else {
			return nil, fmt.Errorf("pares token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*kawaiiMapCliams); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type id invaild")
	}
}

func PassAdminToken(cfg config.IJwtConfig, tokenString string) (*kawaiiMapCliams, error) {
	token, err := jwt.ParseWithClaims(tokenString, &kawaiiMapCliams{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("singing method is invaild")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invaild")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expried")
		} else {
			return nil, fmt.Errorf("pares token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*kawaiiMapCliams); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type id invaild")
	}
}

func RepeatToken(cfg config.IJwtConfig, cliams *users.UserClaims, exp int64) string {
	obj := &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapCliams{
			Claims: cliams,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewKawaiiAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IKawaiiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IKawaiiAuth {
	return &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapCliams{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IKawaiiAuth {
	return &kawaiiAuth{
		cfg: cfg,
		mapClaims: &kawaiiMapCliams{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "kawaiishop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IKawaiiAuth {
	return &kawaiiAdmin{
		kawaiiAuth: &kawaiiAuth{
			cfg: cfg,
			mapClaims: &kawaiiMapCliams{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "kawaiishop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
