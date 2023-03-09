package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Handler struct {
}

type Config struct {
	Engine *gin.Engine
}

func NewHandler(c *Config) {
	h := &Handler{}
	g := c.Engine.Group(os.Getenv("ACCOUNT_API_URL"))
	g.GET("/me", h.Me)
	g.POST("/signup", h.SignUp)
	g.POST("/signin", h.SignIn)
	g.POST("/signout", h.SignOut)
	g.PUT("/details", h.Details)
}

func (h *Handler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it is me page",
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it is sign-up page",
	})
}

func (h *Handler) SignIn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it is sign-in page",
	})
}

func (h *Handler) SignOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it is sign-out page",
	})
}

func (h *Handler) Details(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "it is details page",
	})
}
