package mocks

import (
	"context"
	"github.com/amiosamu/markdown/account/model"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) NewPairFromUser(ctx context.Context, u *model.User, prevToken string) (*model.TokenPair, error) {
	ret := m.Called(ctx, u, prevToken)
	var r0 *model.TokenPair
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.TokenPair)
	}
	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	return r0, r1
}
