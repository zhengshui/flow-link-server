package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhengshui/flow-link-server/api/middleware"
	"github.com/zhengshui/flow-link-server/bootstrap"
	"github.com/zhengshui/flow-link-server/mongo"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, router *gin.Engine) {
	// Health check endpoint (for Docker/K8s health probes)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "flow-link-server",
		})
	})

	// API group with /api prefix
	apiGroup := router.Group("/api")

	// Public APIs (no authentication required)
	publicRouter := apiGroup.Group("")
	NewSignupRouter(env, timeout, db, publicRouter)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)
	// Plan templates public endpoints (GET only)
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
	// Feedback
	NewFeedbackRouter(env, timeout, db, protectedRouter)
	// Plan templates protected endpoints (POST, PUT, DELETE for personal templates)
	NewProtectedPlanTemplateRouter(env, timeout, db, protectedRouter)

	// Admin APIs (JWT authentication + admin role required)
	adminRouter := apiGroup.Group("/admin")
	adminRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	adminRouter.Use(middleware.AdminAuthMiddleware(env.AccessTokenSecret))
	// Admin plan templates management
	NewAdminPlanTemplateRouter(env, timeout, db, adminRouter)
}
