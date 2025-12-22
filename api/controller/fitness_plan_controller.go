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

// CreateFromTemplate godoc
// @Summary      基于模板创建健身计划
// @Description  使用计划模板创建新的健身计划
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreatePlanFromTemplateRequest true "模板创建请求"
// @Success      200 {object} domain.SuccessResponse{data=domain.FitnessPlan} "创建成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/from-template [post]
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

// CreateCustom godoc
// @Summary      创建自定义健身计划
// @Description  创建完全自定义的健身计划
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateCustomPlanRequest true "自定义计划信息"
// @Success      200 {object} domain.SuccessResponse{data=domain.FitnessPlan} "创建成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/custom [post]
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

// GetByID godoc
// @Summary      获取健身计划详情
// @Description  根据ID获取健身计划的详细信息
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Success      200 {object} domain.SuccessResponse{data=domain.FitnessPlan} "获取成功"
// @Failure      400 {object} domain.ErrorResponse "计划ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      404 {object} domain.ErrorResponse "健身计划不存在"
// @Router       /api/plans/{planId} [get]
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

// GetList godoc
// @Summary      获取健身计划列表
// @Description  获取用户的健身计划列表，支持分页和状态筛选
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "页码" default(1)
// @Param        pageSize query int false "每页数量" default(10)
// @Param        status query string false "计划状态"
// @Success      200 {object} domain.SuccessResponse{data=domain.PaginatedData} "获取成功"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans [get]
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

// UpdateStatus godoc
// @Summary      更新计划状态
// @Description  更新健身计划的状态（激活、暂停、完成等）
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Param        request body domain.UpdatePlanStatusRequest true "状态信息"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "更新成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/{planId}/status [put]
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

// CompleteDay godoc
// @Summary      完成训练日
// @Description  标记健身计划中的某一天为已完成
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Param        request body domain.CompleteDayRequest true "完成日期信息"
// @Success      200 {object} domain.SuccessResponse{data=domain.FitnessPlan} "完成成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/{planId}/complete-day [post]
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

// Delete godoc
// @Summary      删除健身计划
// @Description  删除指定ID的健身计划
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "删除成功"
// @Failure      400 {object} domain.ErrorResponse "计划ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/{planId} [delete]
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

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(nil, "删除成功"))
}

// GetProgress godoc
// @Summary      获取计划进度摘要
// @Description  获取健身计划的进度摘要信息
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Success      200 {object} domain.SuccessResponse{data=domain.PlanProgress} "获取成功"
// @Failure      400 {object} domain.ErrorResponse "计划ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      404 {object} domain.ErrorResponse "计划不存在"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/{planId}/progress [get]
func (fc *FitnessPlanController) GetProgress(c *gin.Context) {
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

	progress, err := fc.FitnessPlanUsecase.GetProgress(c, userID, planID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(404, "计划不存在"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(progress))
}

// SkipDay godoc
// @Summary      跳过计划日
// @Description  跳过健身计划中的某一天
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Param        request body domain.SkipDayRequest true "跳过日期信息"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "已跳过"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/{planId}/skip-day [post]
func (fc *FitnessPlanController) SkipDay(c *gin.Context) {
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

	var request domain.SkipDayRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := fc.FitnessPlanUsecase.SkipDay(c, userID, planID, request.DayNumber, request.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(result, "已跳过"))
}

// AdjustDay godoc
// @Summary      临时调整计划日动作
// @Description  临时调整健身计划中某一天的训练动作（仅对该计划实例生效，不修改原模板）
// @Tags         健身计划
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        planId path string true "计划ID"
// @Param        request body domain.AdjustDayRequest true "调整信息"
// @Success      200 {object} domain.SuccessResponse "调整成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/plans/{planId}/adjust-day [post]
func (fc *FitnessPlanController) AdjustDay(c *gin.Context) {
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

	var request domain.AdjustDayRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	err = fc.FitnessPlanUsecase.AdjustDay(c, userID, planID, request.DayNumber, request.Exercises, request.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(nil, "调整成功"))
}
