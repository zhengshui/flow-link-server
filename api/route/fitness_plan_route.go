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

func NewFitnessPlanRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	fp := repository.NewFitnessPlanRepository(db, domain.CollectionFitnessPlan)
	pt := repository.NewPlanTemplateRepository(db, domain.CollectionPlanTemplate)
	fc := &controller.FitnessPlanController{
		FitnessPlanUsecase: usecase.NewFitnessPlanUsecase(fp, pt, timeout),
	}
	group.POST("/plans/from-template", fc.CreateFromTemplate)
	group.POST("/plans/custom", fc.CreateCustom)
	group.GET("/plans/:planId", fc.GetByID)
	group.GET("/plans", fc.GetList)
	group.PUT("/plans/:planId/status", fc.UpdateStatus)
	group.POST("/plans/:planId/complete-day", fc.CompleteDay)
	group.DELETE("/plans/:planId", fc.Delete)
}
