package main

import (
	"fmt"
	"github.com/amiosamu/markdown/account/handler"
	"github.com/amiosamu/markdown/account/repository"
	"github.com/amiosamu/markdown/account/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"time"
)

func inject(d *DataSources) (*gin.Engine, error) {
	log.Println("injecting data sources")

	userRepository := repository.NewUserRepository(d.DB)
	tokenRepository := repository.NewTokenRepository(d.RedisClient)
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})

	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := os.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := os.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key %w", err)
	}

	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")
	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	refreshSecret := os.Getenv("REFRESH_SECRET")
	tokenService := service.NewTokenService(&service.TSConfig{
		TokenRepository:       tokenRepository,
		PrivateKey:            privKey,
		PublicKey:             pubKey,
		RefreshSecret:         refreshSecret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	router := gin.Default()

	baseURL := os.Getenv("ACCOUNT_API_URL")

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}
	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(time.Duration(time.Duration(ht) * time.Second)),
	})
	return router, nil
}
