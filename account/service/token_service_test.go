package service

import (
	"context"
	"fmt"
	"github.com/amiosamu/markdown/account/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestNewPairFromuser(t *testing.T) {
	priv, _ := os.ReadFile("../rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := os.ReadFile("../rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	tokenService := NewTokenService(&TSConfig{
		PrivateKey:    privKey,
		PublicKey:     pubKey,
		RefreshSecret: secret,
	})
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@gmail.com",
		Password: "avalidpassword",
	}
	t.Run("Returns a token pair with proper values", func(t *testing.T) {
		ctx := context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, u, "")
		assert.NoError(t, err)

		var s string
		assert.IsType(t, s, tokenPair.IDToken)
		idTokenClaims := &IDTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.IDToken, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})
		assert.NoError(t, err)

		expectedClaims := []interface{}{
			u.UID,
			u.Email,
			u.Name,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
		}
		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password)

		expiresAt := time.Unix(idTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt := time.Now().Add(15 * time.Minute)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		fmt.Println(refreshTokenClaims.Id)
		assert.IsType(t, s, tokenPair.RefreshToken)
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.StandardClaims.ExpiresAt, 0)
		expectedExpiresAt = time.Now().Add(3 * 24 * time.Hour)

		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})
}
