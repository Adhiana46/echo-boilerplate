package tokenmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/Adhiana46/echo-boilerplate/entity"
	"github.com/Adhiana46/echo-boilerplate/pkg/cache"
	"github.com/golang-jwt/jwt/v4"
)

var (
	cacheBlacklistFmt = "blacklist:%s"

	ErrInvalidToken     = errors.New("invalid token")
	ErrBlacklistedToken = errors.New("token blacklisted")
)

// in-memory blacklisted token
var blacklistedTokens map[string]string = map[string]string{}

type TokenManager struct {
	cache     cache.Cache
	secretKey string
}

func NewTokenManager(secretKey string, cache cache.Cache) *TokenManager {
	return &TokenManager{
		cache:     cache,
		secretKey: secretKey,
	}
}

func (r *TokenManager) GenerateToken(claims *entity.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(r.secretKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Check if token is valid (signature valid, not-expire, not-blacklisted)
func (r *TokenManager) ParseToken(tokenStr string) (*jwt.Token, *entity.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &entity.UserClaims{}, r.secretKeyFn)
	if err != nil {
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

	claims, ok := token.Claims.(entity.UserClaims)
	if !ok {
		return nil, nil, ErrInvalidToken
	}

	return token, &claims, nil
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
	return []byte(r.secretKey), nil
}
