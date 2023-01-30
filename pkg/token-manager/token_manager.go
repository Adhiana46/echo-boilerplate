package tokenmanager

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Adhiana46/echo-boilerplate/config"
	"github.com/Adhiana46/echo-boilerplate/dto"
	"github.com/Adhiana46/echo-boilerplate/pkg/cache"
	"github.com/golang-jwt/jwt/v4"
)

var (
	cacheBlacklistFmt = "blacklist:%s"

	ErrInvalidToken     = errors.New("invalid token")
	ErrBlacklistedToken = errors.New("token blacklisted")
	ErrTokenExpired     = errors.New("token expired")
	ErrIncorrectIssuer  = errors.New("incorrect issuer")
)

// in-memory blacklisted token
var blacklistedTokens map[string]string = map[string]string{}

type TokenManager struct {
	cache     cache.Cache
	secretKey string
	issuer    string
}

func NewTokenManager(jwtCfg *config.JWTConfig, cache cache.Cache) *TokenManager {
	return &TokenManager{
		cache:     cache,
		secretKey: jwtCfg.SecretKey,
		issuer:    jwtCfg.SecretKey,
	}
}

func (r *TokenManager) GenerateToken(claims *dto.UserClaims) (string, error) {
	claims.Issuer = r.issuer
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenStr, err := token.SignedString([]byte(r.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Check if token is valid (signature valid, not-expire, not-blacklisted)
func (r *TokenManager) ParseToken(tokenStr string) (*jwt.Token, *dto.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &dto.UserClaims{}, r.secretKeyFn)
	if err != nil {
		// check if token is expired
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, nil, ErrTokenExpired
		}

		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, ErrInvalidToken
	}

	// check if token is blacklisted
	cacheKey := fmt.Sprintf(cacheBlacklistFmt, tokenStr)
	if r.cache != nil {
		cacheResult, err := r.cache.Get(cacheKey)
		if err != nil && err != cache.ErrCacheNil {
			return nil, nil, err
		}

		if cacheResult == "1" {
			return nil, nil, ErrBlacklistedToken
		}
	} else {
		val, isExists := blacklistedTokens[cacheKey]

		if isExists && val == "1" {
			return nil, nil, ErrBlacklistedToken
		}
	}

	claims, ok := token.Claims.(*dto.UserClaims)
	if !ok {
		log.Println("NOT OK")
		return nil, nil, ErrInvalidToken
	}

	if claims.Issuer != r.issuer {
		return nil, nil, ErrIncorrectIssuer
	}

	return token, claims, nil
}

// blacklist token
func (r *TokenManager) BlacklistToken(tokenStr string) error {
	cacheKey := fmt.Sprintf(cacheBlacklistFmt, tokenStr)

	_, claims, err := r.ParseToken(tokenStr)
	if err != nil {
		return err
	}

	expireAt := claims.ExpiresAt.Time.Unix() - time.Now().Unix()

	if r.cache != nil {
		return r.cache.Set(cacheKey, "1", int32(expireAt))
	} else {
		blacklistedTokens[cacheKey] = "1"
	}

	return nil
}

func (r *TokenManager) secretKeyFn(token *jwt.Token) (interface{}, error) {
	// validate signing algo
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	return []byte(r.secretKey), nil
}
