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

func NewUserInfoRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	uc := &controller.UserInfoController{
		UserInfoUsecase: usecase.NewUserInfoUsecase(ur, timeout),
	}
	group.GET("/user/info", uc.GetUserInfo)
	group.PUT("/user/info", uc.UpdateUserInfo)
}
