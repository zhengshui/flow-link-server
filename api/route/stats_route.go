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

func NewStatsRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTrainingRecordRepository(db, domain.CollectionTrainingRecord)
	fp := repository.NewFitnessPlanRepository(db, domain.CollectionFitnessPlan)
	sc := &controller.StatsController{
		StatsUsecase: usecase.NewStatsUsecase(tr, fp, timeout),
	}
	group.GET("/stats/training", sc.GetTrainingStats)
	group.GET("/stats/muscle-groups", sc.GetMuscleGroupStats)
	group.GET("/stats/personal-records", sc.GetPersonalRecords)
	group.GET("/stats/calendar", sc.GetCalendar)
	// v1.3.0 新增路由
	group.GET("/stats/plan", sc.GetPlanStats)
	group.GET("/stats/plan-progress", sc.GetPlanProgressList)
}
