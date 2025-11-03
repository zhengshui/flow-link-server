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

// GetByID 获取计划模板详情（无需认证）
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

// GetList 获取计划模板列表（无需认证）
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
