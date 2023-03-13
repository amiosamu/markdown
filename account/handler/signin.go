package handler

import (
	"github.com/amiosamu/markdown/account/model"
	"github.com/amiosamu/markdown/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type signIn struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

func (h *Handler) Signin(c *gin.Context) {

	var req signIn
	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}
	ctx := c.Request.Context()
	err := h.UserService.SignIn(ctx, u)
	if err != nil {
		log.Printf("failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")
	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"tokens": tokens,
	})
}
