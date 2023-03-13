package handler

import (
	"fmt"
	"github.com/amiosamu/markdown/account/model/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func bindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())
		err := apperrors.NewUnsupportedMediaType(msg)
		c.JSON(err.StatusCode(), gin.H{
			"error": err,
		})
		return false
	}

	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)
		if errs, ok := err.(validator.ValidationErrors); ok {
			var invalidArgs []invalidArgument
			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					err.Value().(string),
					err.Tag(),
					err.Param(),
				})
			}
			err := apperrors.NewBadRequest("invalid request parameters.")
			c.JSON(err.StatusCode(), gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}
		fallback := apperrors.NewInternalServerError()
		c.JSON(fallback.StatusCode(), gin.H{
			"error": fallback,
		})
	}
	return true
}
