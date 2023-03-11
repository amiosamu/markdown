package service

import (
	"crypto/rsa"
	"github.com/amiosamu/markdown/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"time"
)

func generateIDToken(u *model.User, key *rsa.PrivateKey) (string, error) {
	unixTime := time.Now().Unix()
	tokenExp := unixTime + 60*15

	claims := IDTokenCustomClaims{
		User: u,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime,
			ExpiresAt: tokenExp,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		log.Println("failed to sign id token string")
		return "", err
	}
	return ss, nil
}

type RefreshTokenStruct struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

type RefreshToken struct {
	SS        string
	ID        string
	ExpiresIn time.Duration
}

type IDTokenCustomClaims struct {
	User *model.User `json:"user"`
	jwt.StandardClaims
}

type RefreshTokenCustomClaims struct {
	UID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

func generateRefreshToken(uid uuid.UUID, key string) (*RefreshToken, error) {
	currentTime := time.Now()
	tokenExp := currentTime.AddDate(0, 0, 3)
	tokenID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Failed to generate refresh token ID")
		return nil, err
	}
	claims := RefreshTokenCustomClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  currentTime.Unix(),
			ExpiresAt: tokenExp.Unix(),
			Id:        tokenID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}
	return &RefreshToken{
		SS:        ss,
		ID:        tokenID.String(),
		ExpiresIn: tokenExp.Sub(currentTime),
	}, nil
}