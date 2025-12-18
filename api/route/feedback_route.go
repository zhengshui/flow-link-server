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

func NewFeedbackRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	fr := repository.NewFeedbackRepository(db, domain.CollectionFeedback)
	fc := controller.FeedbackController{
		FeedbackUsecase: usecase.NewFeedbackUsecase(fr, timeout),
	}
	group.POST("/feedback", fc.CreateFeedback)
}
