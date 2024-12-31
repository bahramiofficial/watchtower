package handlers

import (
	"fmt"
	"net/http"

	"github.com/bahramiofficial/watchtower/src/api/helper"
	"github.com/bahramiofficial/watchtower/src/api/model"
	"github.com/bahramiofficial/watchtower/src/api/service"
	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	HuntService *service.HuntService
}

type data struct {
	mobile string `json:"mobile"`
}

func NewTestHandler() *TestHandler {
	return &TestHandler{HuntService: service.NewHuntService()}
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

//1 22.30 s27

func (t *TestHandler) CreateHunt(c *gin.Context) {
	req := model.CreateHuntRequest{}
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseHttpResponseWithError(nil, false, -1, err))
		return
	}
	res, err := t.HuntService.Create(c, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseHttpResponseWithError(nil, false, -1, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseHttpResponse(res, true, 1))
	return
}
