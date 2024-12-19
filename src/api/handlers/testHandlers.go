package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
}

type data struct {
	mobile string `json:"mobile"`
}
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (t *TestHandler) Test(c *gin.Context) {
	c.JSON(http.StatusOK, "working")
	return
}

func (t *TestHandler) TestById(c *gin.Context) {
	c.JSON(http.StatusOK, fmt.Sprintf("id : %s", c.Params.ByName("id")))
	return
}

func (t *TestHandler) TestByDataTestValidationAndBinding(c *gin.Context) {
    da := data{}
	c.ShouldBind(da)

	c.JSON(http.StatusOK, "ok")
	return
}
