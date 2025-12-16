package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FeedbackController struct {
	FeedbackUsecase domain.FeedbackUsecase
}

// CreateFeedback godoc
// @Summary      提交用户反馈
// @Description  用户提交意见反馈
// @Tags         反馈
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.FeedbackRequest true "反馈信息"
// @Success      200 {object} domain.SuccessResponse{data=domain.FeedbackResponse} "反馈提交成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/feedback [post]
func (fc *FeedbackController) CreateFeedback(c *gin.Context) {
	var request domain.FeedbackRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	// 从JWT中获取用户ID
	userID := c.GetString("x-user-id")
	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "无效的用户ID"))
		return
	}

	// 设置默认反馈类型
	feedbackType := request.Type
	if feedbackType == "" {
		feedbackType = "建议"
	}

	// 验证反馈类型
	validTypes := map[string]bool{"建议": true, "问题": true, "其他": true}
	if !validTypes[feedbackType] {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "无效的反馈类型，请选择：建议/问题/其他"))
		return
	}

	// 创建反馈
	now := time.Now()
	feedback := domain.Feedback{
		ID:          primitive.NewObjectID(),
		UserID:      userIDHex,
		Content:     request.Content,
		Type:        feedbackType,
		ContactInfo: request.ContactInfo,
		Status:      "待处理",
	}

	err = fc.FeedbackUsecase.Create(c, &feedback)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	response := domain.FeedbackResponse{
		ID:        feedback.ID.Hex(),
		CreatedAt: now.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, domain.ApiResponse{
		Code:    200,
		Message: "反馈提交成功",
		Data:    response,
	})
}

