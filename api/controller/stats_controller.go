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

// GetTrainingStats godoc
// @Summary      获取训练统计
// @Description  获取用户的训练统计数据，支持按周期查询
// @Tags         统计
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period query string false "统计周期(week/month/year)" default(week)
// @Param        startDate query string false "开始日期" format(date)
// @Param        endDate query string false "结束日期" format(date)
// @Success      200 {object} domain.SuccessResponse{data=domain.TrainingStats} "获取成功"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /stats/training [get]
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

// GetMuscleGroupStats godoc
// @Summary      获取肌群训练统计
// @Description  获取用户的肌群训练统计数据
// @Tags         统计
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period query string false "统计周期(week/month/year)" default(month)
// @Success      200 {object} domain.SuccessResponse{data=[]domain.MuscleGroupStats} "获取成功"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /stats/muscle-groups [get]
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

// GetPersonalRecords godoc
// @Summary      获取个人记录
// @Description  获取用户的个人最佳记录
// @Tags         统计
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.SuccessResponse{data=[]domain.PersonalRecord} "获取成功"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /stats/personal-records [get]
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

// GetCalendar godoc
// @Summary      获取日历数据
// @Description  获取用户指定月份的训练日历数据
// @Tags         统计
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        year query int false "年份" default(2024)
// @Param        month query int false "月份" default(1)
// @Success      200 {object} domain.SuccessResponse{data=[]domain.CalendarDay} "获取成功"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /stats/calendar [get]
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
