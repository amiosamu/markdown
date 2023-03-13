package repository

import (
	"context"
	"fmt"
	"github.com/amiosamu/markdown/account/model"
	"github.com/amiosamu/markdown/account/model/apperrors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redis *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redis,
	}
}

func (r redisTokenRepository) SetRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		log.Printf("could not SET refresh token to redis for userID/tokenID: %s/%s: %v\\n\", userID, tokenID, err")
		return apperrors.NewInternalServerError()
	}
	return nil
}

func (r redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID, tokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)

	if err := r.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("could not DEL refrsesh token to redis for userID/tokenID: %s/%s: %v\n", userID, tokenID, err)
		return apperrors.NewInternalServerError()
	}
	return nil
}
