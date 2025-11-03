package controller

import (
	"net/http"
	"time"

	"github.com/zhengshui/flow-link-server/bootstrap"
	"github.com/zhengshui/flow-link-server/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *bootstrap.Env
}

func (sc *SignupController) Signup(c *gin.Context) {
	var request domain.SignupRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	// 检查用户名是否已存在
	_, err = sc.SignupUsecase.GetUserByUsername(c, request.Username)
	if err == nil {
		c.JSON(http.StatusConflict, domain.NewErrorResponse(409, "用户名已存在"))
		return
	}

	// 加密密码
	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	// 创建用户
	now := time.Now()
	user := domain.User{
		ID:           primitive.NewObjectID(),
		Username:     request.Username,
		Password:     string(encryptedPassword),
		Nickname:     request.Nickname,
		Email:        request.Email,
		Phone:        request.Phone,
		Gender:       request.Gender,
		Age:          request.Age,
		Height:       request.Height,
		Weight:       request.Weight,
		TargetWeight: request.TargetWeight,
		FitnessGoal:  request.FitnessGoal,
		Role:         "user",
		JoinDate:     now.Format("2006-01-02"),
		CreatedAt:    primitive.NewDateTimeFromTime(now),
		UpdatedAt:    primitive.NewDateTimeFromTime(now),
	}

	err = sc.SignupUsecase.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	// 生成token
	accessToken, err := sc.SignupUsecase.CreateAccessToken(&user, sc.Env.AccessTokenSecret, sc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	signupResponse := domain.SignupResponse{
		Token: accessToken,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(signupResponse))
}
