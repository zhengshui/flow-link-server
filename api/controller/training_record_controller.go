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

// Create 创建训练记录
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

// GetByID 获取训练记录详情
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

// GetList 获取训练记录列表
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

// Update 更新训练记录
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

// Delete 删除训练记录
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
