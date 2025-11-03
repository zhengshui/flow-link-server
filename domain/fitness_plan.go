package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionFitnessPlan = "fitness_plans"
)

// TrainingDay 训练日程
type TrainingDay struct {
	DayNumber int        `bson:"dayNumber" json:"dayNumber"` // 第几天
	DayName   string     `bson:"dayName" json:"dayName"`     // 训练日名称
	IsRestDay bool       `bson:"isRestDay" json:"isRestDay"` // 是否为休息日
	Exercises []Exercise `bson:"exercises" json:"exercises"` // 当日训练项目
	Notes     string     `bson:"notes" json:"notes,omitempty"` // 当日备注
}

// FitnessPlan 健身计划
type FitnessPlan struct {
	ID                   primitive.ObjectID `bson:"_id" json:"id"`
	UserID               primitive.ObjectID `bson:"userId" json:"userId"`
	TemplateID           int                `bson:"templateId" json:"templateId"`               // 模板ID(0表示自定义)
	Name                 string             `bson:"name" json:"name"`
	Description          string             `bson:"description" json:"description"`
	Goal                 string             `bson:"goal" json:"goal"`                           // 训练目标
	DurationWeeks        int                `bson:"durationWeeks" json:"durationWeeks"`         // 计划周期(周)
	TrainingDaysPerWeek  int                `bson:"trainingDaysPerWeek" json:"trainingDaysPerWeek"` // 每周训练天数
	TrainingDays         []TrainingDay      `bson:"trainingDays" json:"trainingDays"`           // 训练日程
	StartDate            string             `bson:"startDate" json:"startDate"`                 // 开始日期 YYYY-MM-DD
	EndDate              string             `bson:"endDate" json:"endDate"`                     // 结束日期 YYYY-MM-DD
	Status               string             `bson:"status" json:"status"`                       // 进行中/已完成/已暂停/已归档
	CurrentWeek          int                `bson:"currentWeek" json:"currentWeek"`             // 当前第几周
	CurrentDay           int                `bson:"currentDay" json:"currentDay"`               // 当前第几天
	CompletedDays        []int              `bson:"completedDays" json:"completedDays"`         // 已完成的训练日
	TotalCompletedDays   int                `bson:"totalCompletedDays" json:"totalCompletedDays"` // 累计完成天数
	CompletionRate       int                `bson:"completionRate" json:"completionRate"`       // 完成率(百分比)
	CreatedAt            primitive.DateTime `bson:"createdAt" json:"createdAt" swaggertype:"string"`
	UpdatedAt            primitive.DateTime `bson:"updatedAt" json:"updatedAt" swaggertype:"string"`
}

// FitnessPlanRepository 健身计划仓储接口
type FitnessPlanRepository interface {
	Create(c context.Context, plan *FitnessPlan) error
	GetByID(c context.Context, id string) (FitnessPlan, error)
	GetByUserID(c context.Context, userID string, status string, page, pageSize int) ([]FitnessPlan, int64, error)
	Update(c context.Context, id string, plan *FitnessPlan) error
	UpdateStatus(c context.Context, id string, status string) error
	Delete(c context.Context, id string) error
	CompletePlanDay(c context.Context, id string, dayNumber int) error
}

// CreatePlanFromTemplateRequest 基于模板创建计划请求
type CreatePlanFromTemplateRequest struct {
	TemplateID int    `json:"templateId" binding:"required"`
	StartDate  string `json:"startDate" binding:"required"`
	Name       string `json:"name"`
}

// CreateCustomPlanRequest 创建自定义计划请求
type CreateCustomPlanRequest struct {
	Name                string        `json:"name" binding:"required"`
	Description         string        `json:"description"`
	Goal                string        `json:"goal" binding:"required"`
	DurationWeeks       int           `json:"durationWeeks" binding:"required"`
	TrainingDaysPerWeek int           `json:"trainingDaysPerWeek" binding:"required"`
	TrainingDays        []TrainingDay `json:"trainingDays" binding:"required"`
	StartDate           string        `json:"startDate" binding:"required"`
}

// CompleteDayRequest 标记训练日完成请求
type CompleteDayRequest struct {
	DayNumber int `json:"dayNumber" binding:"required"`
	RecordID  int `json:"recordId"`
}

// UpdatePlanStatusRequest 更新计划状态请求
type UpdatePlanStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// FitnessPlanUsecase 健身计划用例接口
type FitnessPlanUsecase interface {
	CreateFromTemplate(c context.Context, userID string, request *CreatePlanFromTemplateRequest) (map[string]interface{}, error)
	CreateCustom(c context.Context, userID string, request *CreateCustomPlanRequest) (map[string]interface{}, error)
	GetByID(c context.Context, userID, planID string) (FitnessPlan, error)
	GetList(c context.Context, userID string, status string, page, pageSize int) ([]FitnessPlan, int64, error)
	UpdateStatus(c context.Context, userID, planID string, status string) error
	CompleteDay(c context.Context, userID, planID string, dayNumber int) (map[string]interface{}, error)
	Delete(c context.Context, userID, planID string) error
}
