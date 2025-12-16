package controller

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/bootstrap"
	"github.com/zhengshui/flow-link-server/domain"
)

type LoginController struct {
	LoginUsecase domain.LoginUsecase
	Env          *bootstrap.Env
}

// Login godoc
// @Summary      用户登录
// @Description  用户通过用户名和密码登录，返回JWT token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        request body domain.LoginRequest true "登录信息"
// @Success      200 {object} domain.SuccessResponse{data=domain.LoginResponse} "登录成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "用户名或密码错误"
// @Failure      404 {object} domain.ErrorResponse "用户不存在"
// @Router       /api/auth/login [post]
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
