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

func NewPlanTemplateRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	pt := repository.NewPlanTemplateRepository(db, domain.CollectionPlanTemplate)
	pc := &controller.PlanTemplateController{
		PlanTemplateUsecase: usecase.NewPlanTemplateUsecase(pt, timeout),
	}
	group.GET("/templates/:templateId", pc.GetByID)
	group.GET("/templates", pc.GetList)
}
