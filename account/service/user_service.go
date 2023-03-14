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

func (s *userService) SignIn(ctx context.Context, u *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)

	if err != nil {
		return apperrors.NewAuthorization("invalid email and password")
	}
	match, err := comparePasswords(uFetched.Password, u.Password)

	if err != nil {
		return apperrors.NewInternalServerError()
	}

	if !match {
		return apperrors.NewAuthorization("passwords do not match")
	}
	return nil
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

func (s *userService) Signin(ctx context.Context, u *model.User) error {
	// TODO implement

	return nil
}
