package model

import (
	"context"
	"github.com/google/uuid"
)

type UserService interface {
	Get(ctx context.Context, uuid uuid.UUID) (*User, error)
	Signup(ctx context.Context, u *User) error
}

type TokenService interface {
	NewPairFromUser(ctx context.Context, u *User, prevToken string) (*TokenPair, error)
}

type UserRepository interface {
	FindByID(ctx context.Context, uuid uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
}
