package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
)

type UserInfoController struct {
	UserInfoUsecase domain.UserInfoUsecase
}

// GetUserInfo 获取用户信息
func (uc *UserInfoController) GetUserInfo(c *gin.Context) {
	userIDValue, exists := c.Get("x-user-id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "未授权访问"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "用户ID格式错误"))
		return
	}

	user, err := uc.UserInfoUsecase.GetUserInfo(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(404, "用户不存在"))
		return
	}

	// 构建响应数据
	userInfo := domain.UserInfoResponse{
		ID:           0, // MongoDB ObjectID, frontend uses string
		Username:     user.Username,
		Nickname:     user.Nickname,
		AvatarUrl:    user.AvatarUrl,
		Email:        user.Email,
		Phone:        user.Phone,
		Gender:       user.Gender,
		Age:          user.Age,
		Height:       user.Height,
		Weight:       user.Weight,
		TargetWeight: user.TargetWeight,
		FitnessGoal:  user.FitnessGoal,
		JoinDate:     user.JoinDate,
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(userInfo))
}

// UpdateUserInfo 更新用户信息
func (uc *UserInfoController) UpdateUserInfo(c *gin.Context) {
	userIDValue, exists := c.Get("x-user-id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "未授权访问"))
		return
	}

	userID, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.NewErrorResponse(401, "用户ID格式错误"))
		return
	}

	var request domain.UpdateUserInfoRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	err = uc.UserInfoUsecase.UpdateUserInfo(c, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "更新用户信息失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(map[string]interface{}{
		"message": "用户信息更新成功",
	}))
}
