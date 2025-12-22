package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/api/controller"
	"github.com/zhengshui/flow-link-server/bootstrap"
	"github.com/zhengshui/flow-link-server/domain"
	"github.com/zhengshui/flow-link-server/mongo"
	"github.com/zhengshui/flow-link-server/repository"
	"github.com/zhengshui/flow-link-server/usecase"
)

// NewPlanTemplateRouter 公开路由（无需认证）- 获取官方模板列表和详情
func NewPlanTemplateRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pt := repository.NewPlanTemplateRepository(db, domain.CollectionPlanTemplate)
	pc := &controller.PlanTemplateController{
		PlanTemplateUsecase: usecase.NewPlanTemplateUsecase(pt, timeout),
	}
	group.GET("/templates/:templateId", pc.GetByID)
	group.GET("/templates", pc.GetList)
}

// NewProtectedPlanTemplateRouter 受保护路由（需要认证）- 个人模板的创建、更新、删除、复制
func NewProtectedPlanTemplateRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pt := repository.NewPlanTemplateRepository(db, domain.CollectionPlanTemplate)
	pc := &controller.PlanTemplateController{
		PlanTemplateUsecase: usecase.NewPlanTemplateUsecase(pt, timeout),
	}
	group.POST("/templates/custom", pc.CreateCustom)
	group.POST("/templates/:templateId/duplicate", pc.Duplicate)
	group.PUT("/templates/:templateId", pc.Update)
	group.DELETE("/templates/:templateId", pc.Delete)
}

// NewAdminPlanTemplateRouter 管理员路由（需要认证+管理员权限）- 官方模板管理
func NewAdminPlanTemplateRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pt := repository.NewPlanTemplateRepository(db, domain.CollectionPlanTemplate)
	pc := &controller.PlanTemplateController{
		PlanTemplateUsecase: usecase.NewPlanTemplateUsecase(pt, timeout),
	}
	group.POST("/templates", pc.CreateOfficial)
}
