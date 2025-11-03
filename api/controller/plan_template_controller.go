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
// @Router       /plan-templates/{templateId} [get]
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
// @Description  获取计划模板列表，支持分页和筛选（无需认证）
// @Tags         计划模板
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        pageSize query int false "每页数量" default(10)
// @Param        goal query string false "健身目标"
// @Param        level query string false "难度级别"
// @Success      200 {object} domain.SuccessResponse{data=domain.PaginatedData} "获取成功"
// @Failure      500 {object} domain.ErrorResponse "服务器错误"
// @Router       /plan-templates [get]
func (pc *PlanTemplateController) GetList(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析筛选参数
	goal := c.DefaultQuery("goal", "")
	level := c.DefaultQuery("level", "")

	templates, total, err := pc.PlanTemplateUsecase.GetList(c, goal, level, page, pageSize)
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
