package httpgin

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/ads"
	"homework9/internal/app"
	"net/http"
	"strconv"
)

func HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, adrepo.ErrNotAuthor):
		c.JSON(http.StatusForbidden, ErrorResponse(err))
	case errors.Is(err, adrepo.ErrValidate) || errors.Is(err, adrepo.ErrNotCreated) || errors.Is(err, adrepo.ErrWasDeleted):
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

func CreateAd(c *gin.Context, a app.App) {
	var adReq createAdRequest
	if err := c.ShouldBind(&adReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	adResp, err := a.CreateAd(c, adReq.Title, adReq.Text, adReq.UserID)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, AdSuccessResponse(adResp))
}

func ChangeAdStatus(c *gin.Context, a app.App) {
	var adReq changeAdStatusRequest

	strId := c.Param("id")
	adId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("id should be a number")))
		return
	}

	if err := c.ShouldBind(&adReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	adResp, err := a.ChangeAdStatus(c, adId, adReq.UserID, adReq.Published)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, AdSuccessResponse(adResp))
}

func UpdateAd(c *gin.Context, a app.App) {
	var adReq updateAdRequest

	strId := c.Param("id")
	adId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("id should be a number")))
		return
	}

	if err := c.ShouldBind(&adReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	adResp, err := a.UpdateAd(c, adId, adReq.UserID, adReq.Title, adReq.Text)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, AdSuccessResponse(adResp))
}

func ListAds(c *gin.Context, a app.App) {
	var err error
	filter := ads.AdFilter{}
	if filter.Pub, err = strconv.ParseBool(c.Query("pub")); err != nil {
		filter.Pub = true // default: Published = true
	}
	if filter.Auth, err = strconv.ParseInt(c.Query("auth"), 10, 64); err != nil {
		filter.Auth = -1
	}
	filter.Title = c.Query("title")

	adResp, err := a.GetList(c, filter)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, AdListSuccessResponse(adResp))
}

func GetAd(c *gin.Context, a app.App) {
	strId := c.Param("id")
	adId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("id should be a number")))
		return
	}

	adResp, err := a.GetByID(c, adId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, AdSuccessResponse(adResp))
}

func DeleteAd(c *gin.Context, a app.App) {
	strId := c.Param("id")
	adId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("id should be a number")))
		return
	}

	var adReq DeleteAdRequest
	if err := c.ShouldBind(&adReq); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	}

	err = a.DeleteAd(c, adId, adReq.AuthorId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func CreateUser(c *gin.Context, a app.App) {
	var req CreateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	resp, err := a.CreateUser(c, req.Name)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, UserSuccessResponse(resp))
}

func GetUser(c *gin.Context, a app.App) {
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("id should be a number")))
		return
	}

	resp, err := a.GetUser(c, id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, UserSuccessResponse(resp))
}

func DeleteUser(c *gin.Context, a app.App) {
	strId := c.Param("id")
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(fmt.Errorf("id should be a number")))
	}

	err = a.DeleteUser(c, id)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
