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
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	Name                string             `bson:"name" json:"name"`
	Description         string             `bson:"description" json:"description"`
	Goal                string             `bson:"goal" json:"goal"`                           // 训练目标
	Level               string             `bson:"level" json:"level"`                         // 难度等级(初级/中级/高级)
	DurationWeeks       int                `bson:"durationWeeks" json:"durationWeeks"`         // 计划周期(周)
	TrainingDaysPerWeek int                `bson:"trainingDaysPerWeek" json:"trainingDaysPerWeek"` // 每周训练天数
	TrainingDays        []TrainingDay      `bson:"trainingDays" json:"trainingDays"`           // 训练日程
	ImageUrl            string             `bson:"imageUrl" json:"imageUrl,omitempty"`         // 封面图片URL
	Author              string             `bson:"author" json:"author"`                       // 作者/来源
	Tags                []string           `bson:"tags" json:"tags"`                           // 标签
	CreatedAt           primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

// PlanTemplateRepository 计划模板仓储接口
type PlanTemplateRepository interface {
	GetByID(c context.Context, id string) (PlanTemplate, error)
	GetList(c context.Context, goal, level string, page, pageSize int) ([]PlanTemplate, int64, error)
	Create(c context.Context, template *PlanTemplate) error
}

// PlanTemplateUsecase 计划模板用例接口
type PlanTemplateUsecase interface {
	GetByID(c context.Context, templateID string) (PlanTemplate, error)
	GetList(c context.Context, goal, level string, page, pageSize int) ([]PlanTemplate, int64, error)
}
