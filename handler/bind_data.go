package handler

import (
	"github.com/amiosamu/markdown/pkg/apperrors"
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
	}
	return true
}
