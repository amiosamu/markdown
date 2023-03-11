package service

import (
	"context"
	"crypto/rsa"
	"github.com/amiosamu/markdown/model"
	"github.com/amiosamu/markdown/pkg/apperrors"
	"log"
)

type tokenService struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

// tsConfig will hold repositories that will eventually be injected into
// this service layer
type tsConfig struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	RefreshSecret string
}

// NewTokenService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewTokenService(c *tsConfig) model.TokenService {
	return &tokenService{
		PrivateKey:    c.PrivateKey,
		PublicKey:     c.PublicKey,
		RefreshSecret: c.RefreshSecret,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {
	idToken, err := generateIDToken(u, s.PrivateKey)
	if err != nil {
		log.Printf("Error generating idToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternalServerError()
	}
	refreshToken, err := generateRefreshToken(u.UID, s.RefreshSecret)

	if err != nil {
		log.Printf("Error generating refreshToken for uid: %v. Error: %v\n", u.UID, err.Error())
		return nil, apperrors.NewInternalServerError()
	}
	return &model.TokenPair{
		IDToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil

}
