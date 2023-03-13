package service

import (
	"context"
	"github.com/amiosamu/markdown/account/model"
	"github.com/amiosamu/markdown/account/model/apperrors"
	"github.com/google/uuid"
	"log"
)

type userService struct {
	UserRepository model.UserRepository
}

type USConfig struct {
	UserRepository model.UserRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)
	return u, err
}

func (s *userService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("unable to sign up user for email: %v\n", u.Email)
		return apperrors.NewInternalServerError()
	}
	u.Password = pw
	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}
	return nil
}
