package main

import (
	"fmt"
	"github.com/amiosamu/markdown/handler"
	"github.com/amiosamu/markdown/pkg/database"
	"github.com/amiosamu/markdown/pkg/repository"
	"github.com/amiosamu/markdown/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os"
)

func inject(d *database.DataSources) (*gin.Engine, error) {
	log.Println("injecting data sources")

	userRepository := repository.NewUserRepository(d.DB)
	userService := service.NewUserService(&service.USConfig{
		UserRepository: userRepository,
	})

	privKeyFile := os.Getenv("PRIV_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)

	if err != nil {
		return nil, fmt.Errorf("could not read private key: %w", err)
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("PUB_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key: %w", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key %w", err)
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	tokenService := service.NewTokenService(&service.TSConfig{
		PrivateKey:    privKey,
		PublicKey:     pubKey,
		RefreshSecret: refreshSecret,
	})

	router := gin.Default()

	handler.NewHandler(&handler.Config{
		Engine:       router,
		UserService:  userService,
		TokenService: tokenService,
	})
	return router, nil
}
