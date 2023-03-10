package handler

import (
	"github.com/amiosamu/markdown/model"
	"github.com/amiosamu/markdown/pkg/apperrors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (h *Handler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Printf("unable to extract user from reqeust context  %v\n", c)
		err := apperrors.NewInternalServerError()
		c.JSON(err.StatusCode(), gin.H{
			"error": err,
		})
		return
	}
	uid := user.(*model.User).UID

	u, err := h.UserService.Get(c, uid)
	if err != nil {
		log.Printf("unable to find user %v\nv%v", uid, err)
		err := apperrors.NewNotFound("user", uid.String())
		c.JSON(err.StatusCode(), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}
