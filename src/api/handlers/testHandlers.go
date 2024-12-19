package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
}

func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (t *TestHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, "working")
	return
}

func (t *TestHandler) TestById(c *gin.Context) {
	c.JSON(http.StatusOK, "work id")
	return
}
