package route

import (
	"time"

	"github.com/zhengshui/flow-link-server/api/middleware"
	"github.com/zhengshui/flow-link-server/bootstrap"
	"github.com/zhengshui/flow-link-server/mongo"
	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	// API group with /api prefix
	apiGroup := gin.Group("/api")

	// Public APIs (no authentication required)
	publicRouter := apiGroup.Group("")
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)
	// Plan templates are public
	NewPlanTemplateRouter(env, timeout, db, publicRouter)

	// Protected APIs (JWT authentication required)
	protectedRouter := apiGroup.Group("")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	// User info
	NewUserInfoRouter(env, timeout, db, protectedRouter)
	// Training records
	NewTrainingRecordRouter(env, timeout, db, protectedRouter)
	// Fitness plans
	NewFitnessPlanRouter(env, timeout, db, protectedRouter)
	// Stats
	NewStatsRouter(env, timeout, db, protectedRouter)
}
