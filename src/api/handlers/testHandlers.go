package handlers

import (
	"fmt"
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/helper"
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
	c.JSON(http.StatusOK, helper.GenerateBaseHttpResponse("work with base response", true, 0))
	return
}

func (t *TestHandler) TestById(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseHttpResponse(fmt.Sprintf("id : %s", c.Params.ByName("id")), true, 0))
	return
}

func (t *TestHandler) TestByDataTestValidationAndBinding(c *gin.Context) {
	da := data{}
	c.ShouldBind(da)

	c.JSON(http.StatusOK, "ok")
	return
}
