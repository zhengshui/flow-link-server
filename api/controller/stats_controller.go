package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
)

type StatsController struct {
	StatsUsecase domain.StatsUsecase
}

// GetTrainingStats 获取训练统计
func (sc *StatsController) GetTrainingStats(c *gin.Context) {
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

	// 解析查询参数
	period := c.DefaultQuery("period", "week")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	stats, err := sc.StatsUsecase.GetTrainingStats(c, userID, period, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取训练统计失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(stats))
}

// GetMuscleGroupStats 获取肌群训练统计
func (sc *StatsController) GetMuscleGroupStats(c *gin.Context) {
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

	// 解析查询参数
	period := c.DefaultQuery("period", "month")

	stats, err := sc.StatsUsecase.GetMuscleGroupStats(c, userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取肌群统计失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(stats))
}

// GetPersonalRecords 获取个人记录
func (sc *StatsController) GetPersonalRecords(c *gin.Context) {
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

	records, err := sc.StatsUsecase.GetPersonalRecords(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取个人记录失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(records))
}

// GetCalendar 获取日历数据
func (sc *StatsController) GetCalendar(c *gin.Context) {
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

	// 解析查询参数
	year, _ := strconv.Atoi(c.DefaultQuery("year", "2024"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "1"))

	calendar, err := sc.StatsUsecase.GetCalendar(c, userID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取日历数据失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(calendar))
}
