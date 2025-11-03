package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/zhengshui/flow-link-server/bootstrap"
	"github.com/zhengshui/flow-link-server/domain"
	"github.com/gin-gonic/gin"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

func (lc *LoginController) Login(c *gin.Context) {
	var request domain.LoginRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	user, err := lc.LoginUsecase.GetUserByUsername(c, request.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(404, "用户名或密码错误"))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "用户名或密码错误"))
		return
	}

	accessToken, err := lc.LoginUsecase.CreateAccessToken(&user, lc.Env.AccessTokenSecret, lc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	loginResponse := domain.LoginResponse{
		Token: accessToken,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(loginResponse))
}
