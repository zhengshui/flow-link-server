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
	DayNumber    int        `bson:"dayNumber" json:"dayNumber"`                         // 第几天
	DayName      string     `bson:"dayName" json:"dayName"`                             // 训练日名称
	IsRestDay    bool       `bson:"isRestDay" json:"isRestDay"`                         // 是否为休息日
	Exercises    []Exercise `bson:"exercises" json:"exercises"`                         // 当日训练项目
	Notes        string     `bson:"notes,omitempty" json:"notes,omitempty"`             // 当日备注
	IntensityHint string    `bson:"intensityHint,omitempty" json:"intensityHint,omitempty"` // 强度提示(如百分比/RPE)
	WarmupTips   string     `bson:"warmupTips,omitempty" json:"warmupTips,omitempty"`   // 热身建议
	CooldownTips string     `bson:"cooldownTips,omitempty" json:"cooldownTips,omitempty"` // 放松建议
}

// FitnessPlan 健身计划
type FitnessPlan struct {
	ID                    primitive.ObjectID  `bson:"_id" json:"id"`
	UserID                primitive.ObjectID  `bson:"userId" json:"userId"`
	TemplateID            *primitive.ObjectID `bson:"templateId,omitempty" json:"templateId,omitempty"` // 模板ID(null表示自定义)
	Name                  string              `bson:"name" json:"name"`
	Description           string              `bson:"description" json:"description"`
	Goal                  string              `bson:"goal" json:"goal"`                                     // 训练目标
	DurationWeeks         int                 `bson:"durationWeeks" json:"durationWeeks"`                   // 计划周期(周)
	DurationWeeksOverride *int                `bson:"durationWeeksOverride,omitempty" json:"durationWeeksOverride,omitempty"` // 可选，覆盖模板周期
	TrainingDaysPerWeek   int                 `bson:"trainingDaysPerWeek" json:"trainingDaysPerWeek"`       // 每周训练天数
	TrainingDays          []TrainingDay       `bson:"trainingDays" json:"trainingDays"`                     // 训练日程
	TrainingDaysOverride  []TrainingDay       `bson:"trainingDaysOverride,omitempty" json:"trainingDaysOverride,omitempty"` // 可选，覆盖后的日程
	StartDate             string              `bson:"startDate" json:"startDate"`                           // 开始日期 YYYY-MM-DD
	EndDate               string              `bson:"endDate" json:"endDate"`                               // 结束日期 YYYY-MM-DD
	Status                string              `bson:"status" json:"status"`                                 // 进行中/已完成/已暂停/已归档
	CurrentWeek           int                 `bson:"currentWeek" json:"currentWeek"`                       // 当前第几周
	CurrentDay            int                 `bson:"currentDay" json:"currentDay"`                         // 当前第几天
	CompletedDays         []int               `bson:"completedDays" json:"completedDays"`                   // 已完成的训练日
	SkippedDays           []int               `bson:"skippedDays" json:"skippedDays"`                       // 跳过的训练日
	TotalCompletedDays    int                 `bson:"totalCompletedDays" json:"totalCompletedDays"`         // 累计完成天数
	CompletionRate        int                 `bson:"completionRate" json:"completionRate"`                 // 完成率(百分比)
	TotalWeight           float64             `bson:"totalWeight" json:"totalWeight"`                       // 计划累计重量
	TotalDuration         int                 `bson:"totalDuration" json:"totalDuration"`                   // 计划累计时长
	TotalCalories         int                 `bson:"totalCalories" json:"totalCalories"`                   // 计划累计消耗卡路里
	CreatedAt             primitive.DateTime  `bson:"createdAt" json:"createdAt" swaggertype:"string"`
	UpdatedAt             primitive.DateTime  `bson:"updatedAt" json:"updatedAt" swaggertype:"string"`
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
	SkipPlanDay(c context.Context, id string, dayNumber int) error
	UpdateTrainingDay(c context.Context, id string, dayNumber int, exercises []Exercise, notes string) error
}

// CreatePlanFromTemplateRequest 基于模板创建计划请求
type CreatePlanFromTemplateRequest struct {
	TemplateID            string        `json:"templateId" binding:"required"` // 改为string类型，使用ObjectID
	StartDate             string        `json:"startDate" binding:"required"`
	Name                  string        `json:"name"`
	DurationWeeksOverride *int          `json:"durationWeeksOverride,omitempty"` // 可选，覆盖模板周期
	TrainingDaysOverride  []TrainingDay `json:"trainingDaysOverride,omitempty"`  // 可选，轻量调整日程
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

// SkipDayRequest 跳过计划日请求
type SkipDayRequest struct {
	DayNumber int    `json:"dayNumber" binding:"required"`
	Reason    string `json:"reason"`
}

// AdjustDayRequest 临时调整计划日动作请求
type AdjustDayRequest struct {
	DayNumber int        `json:"dayNumber" binding:"required"`
	Exercises []Exercise `json:"exercises" binding:"required"`
	Notes     string     `json:"notes"`
}

// UpdatePlanStatusRequest 更新计划状态请求
type UpdatePlanStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// PlanProgress 计划进度摘要
type PlanProgress struct {
	PlanID           string  `json:"planId"`
	TotalDays        int     `json:"totalDays"`
	CompletedDays    int     `json:"completedDays"`
	SkippedDays      int     `json:"skippedDays"`
	CompletionRate   int     `json:"completionRate"`
	CurrentWeek      int     `json:"currentWeek"`
	CurrentDay       int     `json:"currentDay"`
	NextTrainingDate string  `json:"nextTrainingDate"`
	TotalDuration    int     `json:"totalDuration"`
	TotalWeight      float64 `json:"totalWeight"`
	TotalCalories    int     `json:"totalCalories"`
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
	GetProgress(c context.Context, userID, planID string) (PlanProgress, error)
	SkipDay(c context.Context, userID, planID string, dayNumber int, reason string) (map[string]interface{}, error)
	AdjustDay(c context.Context, userID, planID string, dayNumber int, exercises []Exercise, notes string) error
}
