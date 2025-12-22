package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionPlanTemplate = "plan_templates"
)

// PlanTemplate 计划模板
type PlanTemplate struct {
	ID                   primitive.ObjectID  `bson:"_id" json:"id"`
	UserID               *primitive.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"` // 用户ID(个人模板才有)
	Name                 string              `bson:"name" json:"name"`
	Description          string              `bson:"description" json:"description"`
	Goal                 string              `bson:"goal" json:"goal"`                                                     // 训练目标
	SplitType            string              `bson:"splitType,omitempty" json:"splitType,omitempty"`                       // 分化方式(二分化/三分化/推拉腿/上下肢/四分化/五分化)
	Level                string              `bson:"level" json:"level"`                                                   // 难度等级(初级/中级/高级)
	Equipment            string              `bson:"equipment,omitempty" json:"equipment,omitempty"`                       // 主要器械(徒手/哑铃/器械/混合)
	DurationWeeks        int                 `bson:"durationWeeks" json:"durationWeeks"`                                   // 计划周期(周)
	TrainingDaysPerWeek  int                 `bson:"trainingDaysPerWeek" json:"trainingDaysPerWeek"`                       // 每周训练天数
	TrainingDays         []TrainingDay       `bson:"trainingDays" json:"trainingDays"`                                     // 训练日程
	ImageUrl             string              `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`                         // 封面图片URL
	Author               string              `bson:"author" json:"author"`                                                 // 作者/来源
	Tags                 []string            `bson:"tags" json:"tags"`                                                     // 标签
	RecommendedIntensity string              `bson:"recommendedIntensity,omitempty" json:"recommendedIntensity,omitempty"` // 推荐强度(如RPE 7-8)
	IsOfficial           bool                `bson:"isOfficial" json:"isOfficial"`                                         // 是否为官方模板
	CreatedAt            primitive.DateTime  `bson:"createdAt" json:"createdAt" swaggertype:"string"`
	UpdatedAt            primitive.DateTime  `bson:"updatedAt,omitempty" json:"updatedAt,omitempty" swaggertype:"string"`
}

// CreateCustomTemplateRequest 创建个人模板请求
type CreateCustomTemplateRequest struct {
	Name                 string        `json:"name" binding:"required"`
	Description          string        `json:"description"`
	Goal                 string        `json:"goal" binding:"required"`
	SplitType            string        `json:"splitType"` // 分化方式
	Level                string        `json:"level"`     // 初级/中级/高级
	Equipment            string        `json:"equipment"` // 徒手/哑铃/器械/混合
	DurationWeeks        int           `json:"durationWeeks" binding:"required"`
	TrainingDaysPerWeek  int           `json:"trainingDaysPerWeek" binding:"required"`
	TrainingDays         []TrainingDay `json:"trainingDays" binding:"required"`
	Tags                 []string      `json:"tags"`
	ImageUrl             string        `json:"imageUrl"`
	RecommendedIntensity string        `json:"recommendedIntensity"`
}

// UpdateTemplateRequest 更新模板请求
type UpdateTemplateRequest struct {
	Name                 *string       `json:"name,omitempty"`
	Description          *string       `json:"description,omitempty"`
	Goal                 *string       `json:"goal,omitempty"`
	SplitType            *string       `json:"splitType,omitempty"`
	Level                *string       `json:"level,omitempty"`
	Equipment            *string       `json:"equipment,omitempty"`
	DurationWeeks        *int          `json:"durationWeeks,omitempty"`
	TrainingDaysPerWeek  *int          `json:"trainingDaysPerWeek,omitempty"`
	TrainingDays         []TrainingDay `json:"trainingDays,omitempty"`
	Tags                 []string      `json:"tags,omitempty"`
	ImageUrl             *string       `json:"imageUrl,omitempty"`
	RecommendedIntensity *string       `json:"recommendedIntensity,omitempty"`
}

// CreateOfficialTemplateRequest 创建官方模板请求（管理员）
type CreateOfficialTemplateRequest struct {
	Name                 string        `json:"name" binding:"required"`
	Description          string        `json:"description"`
	Goal                 string        `json:"goal" binding:"required"`
	SplitType            string        `json:"splitType"` // 全身训练/二分化/三分化/推拉腿/上下肢/四分化/五分化
	Level                string        `json:"level"`     // 初级/中级/高级
	Equipment            string        `json:"equipment"` // 徒手/哑铃/器械/混合
	DurationWeeks        int           `json:"durationWeeks" binding:"required"`
	TrainingDaysPerWeek  int           `json:"trainingDaysPerWeek" binding:"required"`
	TrainingDays         []TrainingDay `json:"trainingDays" binding:"required"`
	Author               string        `json:"author"` // 作者/来源
	Tags                 []string      `json:"tags"`
	RecommendedIntensity string        `json:"recommendedIntensity"` // 推荐强度
	ImageUrl             string        `json:"imageUrl"`
}

// PlanTemplateRepository 计划模板仓储接口
type PlanTemplateRepository interface {
	GetByID(c context.Context, id string) (PlanTemplate, error)
	GetList(c context.Context, goal, level, splitType, equipment string, durationWeeksMin, durationWeeksMax, page, pageSize int) ([]PlanTemplate, int64, error)
	GetUserTemplates(c context.Context, userID string, page, pageSize int) ([]PlanTemplate, int64, error)
	Create(c context.Context, template *PlanTemplate) error
	Update(c context.Context, id string, template *PlanTemplate) error
	Delete(c context.Context, id string) error
}

// PlanTemplateUsecase 计划模板用例接口
type PlanTemplateUsecase interface {
	GetByID(c context.Context, templateID string) (PlanTemplate, error)
	GetList(c context.Context, goal, level, splitType, equipment string, durationWeeksMin, durationWeeksMax, page, pageSize int) ([]PlanTemplate, int64, error)
	CreateCustom(c context.Context, userID string, request *CreateCustomTemplateRequest) (map[string]interface{}, error)
	CreateOfficial(c context.Context, request *CreateOfficialTemplateRequest) (map[string]interface{}, error)
	Duplicate(c context.Context, userID, templateID string) (map[string]interface{}, error)
	Update(c context.Context, userID, templateID string, request *UpdateTemplateRequest) error
	Delete(c context.Context, userID, templateID string) error
}
