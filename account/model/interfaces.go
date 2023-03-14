package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type UserService interface {
	Get(ctx context.Context, uuid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
	SignIn(ctx context.Context, u *User) error
}

type TokenService interface {
	NewPairFromUser(ctx context.Context, u *User, prevToken string) (*TokenPair, error)
}

type UserRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID, prevTokenID string) error
}
