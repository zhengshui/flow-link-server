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

func NewTrainingRecordRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	tr := repository.NewTrainingRecordRepository(db, domain.CollectionTrainingRecord)
	tc := &controller.TrainingRecordController{
		TrainingRecordUsecase: usecase.NewTrainingRecordUsecase(tr, timeout),
	}
	group.POST("/training/records", tc.Create)
	group.GET("/training/records/:recordId", tc.GetByID)
	group.GET("/training/records", tc.GetList)
	group.PUT("/training/records/:recordId", tc.Update)
	group.DELETE("/training/records/:recordId", tc.Delete)
}
