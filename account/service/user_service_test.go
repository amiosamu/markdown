package service

import (
	"context"
	"fmt"
	"github.com/amiosamu/markdown/account/model"
	"github.com/amiosamu/markdown/account/model/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})
		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserService(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(nil, fmt.Errorf("Some error down the call chain"))

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.Nil(t, u)
		assert.Error(t, err)
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignIn(t *testing.T) {
	email := "bob@gmail.com"
	validPassw := "validpassword"
	hashedValidPass, _ := hashPassword(validPassw)
	invalidPass := "invalidpassword"

	mockUserRepository := new(mocks.MockUserRepository)

	us := NewUserService(&USConfig{
		UserRepository: mockUserRepository,
	})

	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    email,
			Password: invalidPass,
		}
		mockUserResp := &model.User{
			UID:      uid,
			Email:    email,
			Password: hashedValidPass,
		}
		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			email,
		}
		mockUserRepository.
			On("FindByEmail", mockArgs...).Return(mockUserResp, nil)

		ctx := context.TODO()
		err := us.SignIn(ctx, mockUser)
		assert.NoError(t, err)
		mockUserRepository.AssertCalled(t, "FindByEmail", mockArgs...)
	})
}
