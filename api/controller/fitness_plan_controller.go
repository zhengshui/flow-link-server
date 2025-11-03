package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
)

type FitnessPlanController struct {
	FitnessPlanUsecase domain.FitnessPlanUsecase
}

// CreateFromTemplate 基于模板创建计划
func (fc *FitnessPlanController) CreateFromTemplate(c *gin.Context) {
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

	var request domain.CreatePlanFromTemplateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := fc.FitnessPlanUsecase.CreateFromTemplate(c, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "创建健身计划失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(result))
}

// CreateCustom 创建自定义计划
func (fc *FitnessPlanController) CreateCustom(c *gin.Context) {
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

	var request domain.CreateCustomPlanRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := fc.FitnessPlanUsecase.CreateCustom(c, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "创建自定义计划失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(result))
}

// GetByID 获取健身计划详情
func (fc *FitnessPlanController) GetByID(c *gin.Context) {
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

	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "计划ID不能为空"))
		return
	}

	plan, err := fc.FitnessPlanUsecase.GetByID(c, userID, planID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(404, "健身计划不存在"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(plan))
}

// GetList 获取健身计划列表
func (fc *FitnessPlanController) GetList(c *gin.Context) {
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

	// 解析状态参数
	status := c.DefaultQuery("status", "")

	plans, total, err := fc.FitnessPlanUsecase.GetList(c, userID, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取健身计划列表失败"))
		return
	}

	paginatedData := domain.PaginatedData{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Plans:    plans,
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(paginatedData))
}

// UpdateStatus 更新计划状态
func (fc *FitnessPlanController) UpdateStatus(c *gin.Context) {
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

	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "计划ID不能为空"))
		return
	}

	var request domain.UpdatePlanStatusRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	err = fc.FitnessPlanUsecase.UpdateStatus(c, userID, planID, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "更新计划状态失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(map[string]interface{}{
		"message": "计划状态更新成功",
	}))
}

// CompleteDay 完成训练日
func (fc *FitnessPlanController) CompleteDay(c *gin.Context) {
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

	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "计划ID不能为空"))
		return
	}

	var request domain.CompleteDayRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := fc.FitnessPlanUsecase.CompleteDay(c, userID, planID, request.DayNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "完成训练日失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(result))
}

// Delete 删除健身计划
func (fc *FitnessPlanController) Delete(c *gin.Context) {
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

	planID := c.Param("planId")
	if planID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "计划ID不能为空"))
		return
	}

	err := fc.FitnessPlanUsecase.Delete(c, userID, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "删除健身计划失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(map[string]interface{}{
		"message": "健身计划删除成功",
	}))
}
