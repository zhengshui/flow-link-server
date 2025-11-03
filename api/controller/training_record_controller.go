package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
)

type TrainingRecordController struct {
	TrainingRecordUsecase domain.TrainingRecordUsecase
}

// Create godoc
// @Summary      创建训练记录
// @Description  创建新的训练记录
// @Tags         训练记录
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateTrainingRecordRequest true "训练记录信息"
// @Success      200 {object} domain.SuccessResponse{data=domain.TrainingRecord} "创建成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /training-records [post]
func (tc *TrainingRecordController) Create(c *gin.Context) {
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

	var request domain.CreateTrainingRecordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := tc.TrainingRecordUsecase.Create(c, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "创建训练记录失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(result))
}

// GetByID godoc
// @Summary      获取训练记录详情
// @Description  根据ID获取训练记录详细信息
// @Tags         训练记录
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        recordId path string true "记录ID"
// @Success      200 {object} domain.SuccessResponse{data=domain.TrainingRecord} "获取成功"
// @Failure      400 {object} domain.ErrorResponse "记录ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      404 {object} domain.ErrorResponse "训练记录不存在"
// @Router       /training-records/{recordId} [get]
func (tc *TrainingRecordController) GetByID(c *gin.Context) {
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

	recordID := c.Param("recordId")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "记录ID不能为空"))
		return
	}

	record, err := tc.TrainingRecordUsecase.GetByID(c, userID, recordID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(404, "训练记录不存在"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(record))
}

// GetList godoc
// @Summary      获取训练记录列表
// @Description  获取训练记录列表，支持分页和筛选
// @Tags         训练记录
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1)
// @Param        pageSize query int false "每页数量" default(10)
// @Param        startDate query string false "开始日期" format(date)
// @Param        endDate query string false "结束日期" format(date)
// @Param        planId query int false "计划ID"
// @Success      200 {object} domain.SuccessResponse{data=domain.PaginatedData} "获取成功"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /training-records [get]
func (tc *TrainingRecordController) GetList(c *gin.Context) {
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

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析日期范围参数
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	// 解析计划ID参数
	planID, _ := strconv.Atoi(c.DefaultQuery("planId", "0"))

	records, total, err := tc.TrainingRecordUsecase.GetList(c, userID, page, pageSize, startDate, endDate, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取训练记录列表失败"))
		return
	}

	paginatedData := domain.PaginatedData{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Records:  records,
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(paginatedData))
}

// Update godoc
// @Summary      更新训练记录
// @Description  更新指定ID的训练记录
// @Tags         训练记录
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        recordId path string true "记录ID"
// @Param        request body domain.UpdateTrainingRecordRequest true "训练记录信息"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "更新成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /training-records/{recordId} [put]
func (tc *TrainingRecordController) Update(c *gin.Context) {
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

	recordID := c.Param("recordId")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "记录ID不能为空"))
		return
	}

	var request domain.UpdateTrainingRecordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	err = tc.TrainingRecordUsecase.Update(c, userID, recordID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "更新训练记录失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(map[string]interface{}{
		"message": "训练记录更新成功",
	}))
}

// Delete godoc
// @Summary      删除训练记录
// @Description  删除指定ID的训练记录
// @Tags         训练记录
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        recordId path string true "记录ID"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "删除成功"
// @Failure      400 {object} domain.ErrorResponse "记录ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /training-records/{recordId} [delete]
func (tc *TrainingRecordController) Delete(c *gin.Context) {
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

	recordID := c.Param("recordId")
	if recordID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "记录ID不能为空"))
		return
	}

	err := tc.TrainingRecordUsecase.Delete(c, userID, recordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "删除训练记录失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(map[string]interface{}{
		"message": "训练记录删除成功",
	}))
}
