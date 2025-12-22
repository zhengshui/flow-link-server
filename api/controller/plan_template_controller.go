package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/domain"
)

type PlanTemplateController struct {
	PlanTemplateUsecase domain.PlanTemplateUsecase
}

// GetByID godoc
// @Summary      获取计划模板详情
// @Description  根据ID获取计划模板的详细信息（无需认证）
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Param        templateId path string true "模板ID"
// @Success      200 {object} domain.SuccessResponse{data=domain.PlanTemplate} "获取成功"
// @Failure      400 {object} domain.ErrorResponse "模板ID不能为空"
// @Failure      404 {object} domain.ErrorResponse "计划模板不存在"
// @Router       /api/templates/{templateId} [get]
func (pc *PlanTemplateController) GetByID(c *gin.Context) {
	templateID := c.Param("templateId")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "模板ID不能为空"))
		return
	}

	template, err := pc.PlanTemplateUsecase.GetByID(c, templateID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.NewErrorResponse(404, "计划模板不存在"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(template))
}

// GetList godoc
// @Summary      获取计划模板列表
// @Description  获取官方计划模板列表，支持分页和筛选（无需认证）
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        pageSize query int false "每页数量" default(20)
// @Param        goal query string false "训练目标(增肌/减脂/力量提升/耐力提升/综合健身)"
// @Param        level query string false "难度等级(初级/中级/高级)"
// @Param        splitType query string false "分化方式(二分化/三分化/推拉腿/上下肢/四分化/五分化)"
// @Param        equipment query string false "主要器械(徒手/哑铃/器械/混合)"
// @Param        durationWeeksMin query int false "最小周期(周)"
// @Param        durationWeeksMax query int false "最长周期(周)"
// @Success      200 {object} domain.SuccessResponse{data=domain.PaginatedData} "获取成功"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/templates [get]
func (pc *PlanTemplateController) GetList(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// 解析筛选参数
	goal := c.DefaultQuery("goal", "")
	level := c.DefaultQuery("level", "")
	splitType := c.DefaultQuery("splitType", "")
	equipment := c.DefaultQuery("equipment", "")
	durationWeeksMin, _ := strconv.Atoi(c.DefaultQuery("durationWeeksMin", "0"))
	durationWeeksMax, _ := strconv.Atoi(c.DefaultQuery("durationWeeksMax", "0"))

	templates, total, err := pc.PlanTemplateUsecase.GetList(c, goal, level, splitType, equipment, durationWeeksMin, durationWeeksMax, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "获取计划模板列表失败"))
		return
	}

	paginatedData := domain.PaginatedData{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Templates: templates,
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponse(paginatedData))
}

// CreateOfficial godoc
// @Summary      创建官方模板（管理员）
// @Description  管理员创建官方计划模板
// @Tags         管理员-计划模板
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateOfficialTemplateRequest true "模板信息"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "创建成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      403 {object} domain.ErrorResponse "需要管理员权限"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/admin/templates [post]
func (pc *PlanTemplateController) CreateOfficial(c *gin.Context) {
	var request domain.CreateOfficialTemplateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := pc.PlanTemplateUsecase.CreateOfficial(c, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "创建官方模板失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(result, "创建成功"))
}

// CreateCustom godoc
// @Summary      创建个人模板
// @Description  创建个人计划模板
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateCustomTemplateRequest true "模板信息"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "创建成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/templates/custom [post]
func (pc *PlanTemplateController) CreateCustom(c *gin.Context) {
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

	var request domain.CreateCustomTemplateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	result, err := pc.PlanTemplateUsecase.CreateCustom(c, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "创建个人模板失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(result, "创建成功"))
}

// Duplicate godoc
// @Summary      复制官方模板为个人模板
// @Description  复制官方模板创建为用户的个人模板
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        templateId path string true "模板ID"
// @Success      200 {object} domain.SuccessResponse{data=map[string]interface{}} "复制成功"
// @Failure      400 {object} domain.ErrorResponse "模板ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      404 {object} domain.ErrorResponse "模板不存在"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/templates/{templateId}/duplicate [post]
func (pc *PlanTemplateController) Duplicate(c *gin.Context) {
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

	templateID := c.Param("templateId")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "模板ID不能为空"))
		return
	}

	result, err := pc.PlanTemplateUsecase.Duplicate(c, userID, templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(result, "复制成功"))
}

// Update godoc
// @Summary      更新个人模板
// @Description  更新用户的个人计划模板（仅模板所有者）
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        templateId path string true "模板ID"
// @Param        request body domain.UpdateTemplateRequest true "模板信息"
// @Success      200 {object} domain.SuccessResponse "更新成功"
// @Failure      400 {object} domain.ErrorResponse "请求参数错误"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      403 {object} domain.ErrorResponse "无权限修改"
// @Failure      404 {object} domain.ErrorResponse "模板不存在"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/templates/{templateId} [put]
func (pc *PlanTemplateController) Update(c *gin.Context) {
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

	templateID := c.Param("templateId")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "模板ID不能为空"))
		return
	}

	var request domain.UpdateTemplateRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, err.Error()))
		return
	}

	err = pc.PlanTemplateUsecase.Update(c, userID, templateID, &request)
	if err != nil {
		if err.Error() == "unauthorized: you can only update your own templates" {
			c.JSON(http.StatusForbidden, domain.NewErrorResponse(403, "无权限修改该模板"))
			return
		}
		if err.Error() == "cannot modify official templates" {
			c.JSON(http.StatusForbidden, domain.NewErrorResponse(403, "不能修改官方模板"))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "更新模板失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(nil, "更新成功"))
}

// Delete godoc
// @Summary      删除个人模板
// @Description  删除用户的个人计划模板（仅模板所有者）
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        templateId path string true "模板ID"
// @Success      200 {object} domain.SuccessResponse "删除成功"
// @Failure      400 {object} domain.ErrorResponse "模板ID不能为空"
// @Failure      401 {object} domain.ErrorResponse "未授权访问"
// @Failure      403 {object} domain.ErrorResponse "无权限删除"
// @Failure      404 {object} domain.ErrorResponse "模板不存在"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /api/templates/{templateId} [delete]
func (pc *PlanTemplateController) Delete(c *gin.Context) {
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

	templateID := c.Param("templateId")
	if templateID == "" {
		c.JSON(http.StatusBadRequest, domain.NewErrorResponse(400, "模板ID不能为空"))
		return
	}

	err := pc.PlanTemplateUsecase.Delete(c, userID, templateID)
	if err != nil {
		if err.Error() == "unauthorized: you can only delete your own templates" {
			c.JSON(http.StatusForbidden, domain.NewErrorResponse(403, "无权限删除该模板"))
			return
		}
		if err.Error() == "cannot delete official templates" {
			c.JSON(http.StatusForbidden, domain.NewErrorResponse(403, "不能删除官方模板"))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(500, "删除模板失败"))
		return
	}

	c.JSON(http.StatusOK, domain.NewSuccessResponseWithMessage(nil, "删除成功"))
}
