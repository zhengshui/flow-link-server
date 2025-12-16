package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionTrainingRecord = "training_records"
)

// Exercise 训练项目
type Exercise struct {
	ID          int      `bson:"id" json:"id"`
	Name        string   `bson:"name" json:"name"`                                   // 项目名称(必填)
	Sets        *int     `bson:"sets,omitempty" json:"sets,omitempty"`               // 组数
	Reps        *int     `bson:"reps,omitempty" json:"reps,omitempty"`               // 次数
	Weight      *float64 `bson:"weight,omitempty" json:"weight,omitempty"`           // 重量(kg)
	RestTime    *int     `bson:"restTime,omitempty" json:"restTime,omitempty"`       // 休息时间(秒)
	MuscleGroup *string  `bson:"muscleGroup,omitempty" json:"muscleGroup,omitempty"` // 目标肌群
	Notes       *string  `bson:"notes,omitempty" json:"notes,omitempty"`             // 备注
	Duration    *int     `bson:"duration,omitempty" json:"duration,omitempty"`       // 训练时长(分钟)
}

// TrainingRecord 训练记录
type TrainingRecord struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Title          string             `bson:"title" json:"title"`                                       // 标题(必填)
	StartTime      *string            `bson:"startTime,omitempty" json:"startTime,omitempty"`           // 开始时间 YYYY-MM-DD HH:mm:ss
	EndTime        *string            `bson:"endTime,omitempty" json:"endTime,omitempty"`               // 结束时间 YYYY-MM-DD HH:mm:ss
	Duration       *int               `bson:"duration,omitempty" json:"duration,omitempty"`             // 总时长(分钟)
	Exercises      []Exercise         `bson:"exercises,omitempty" json:"exercises,omitempty"`           // 训练项目列表
	TotalWeight    *float64           `bson:"totalWeight,omitempty" json:"totalWeight,omitempty"`       // 总重量(kg)
	TotalSets      *int               `bson:"totalSets,omitempty" json:"totalSets,omitempty"`           // 总组数
	CaloriesBurned *int               `bson:"caloriesBurned,omitempty" json:"caloriesBurned,omitempty"` // 消耗卡路里
	Notes          *string            `bson:"notes,omitempty" json:"notes,omitempty"`                   // 训练备注
	Mood           *string            `bson:"mood,omitempty" json:"mood,omitempty"`                     // 训练状态(优秀/良好/一般/疲劳)
	PlanID         *int               `bson:"planId,omitempty" json:"planId,omitempty"`                 // 关联计划ID
	CreatedAt      primitive.DateTime `bson:"createdAt" json:"createdAt" swaggertype:"string"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt" json:"updatedAt" swaggertype:"string"`
}

// TrainingRecordRepository 训练记录仓储接口
type TrainingRecordRepository interface {
	Create(c context.Context, record *TrainingRecord) error
	GetByID(c context.Context, id string) (TrainingRecord, error)
	GetByUserID(c context.Context, userID string, page, pageSize int, startDate, endDate string, planID int) ([]TrainingRecord, int64, error)
	Update(c context.Context, id string, record *TrainingRecord) error
	Delete(c context.Context, id string) error
}

// CreateTrainingRecordRequest 创建训练记录请求
type CreateTrainingRecordRequest struct {
	Title          string     `json:"title" binding:"required"` // 标题(必填)
	StartTime      *string    `json:"startTime,omitempty"`      // 开始时间 YYYY-MM-DD HH:mm:ss
	EndTime        *string    `json:"endTime,omitempty"`        // 结束时间 YYYY-MM-DD HH:mm:ss
	Duration       *int       `json:"duration,omitempty"`       // 总时长(分钟)
	Exercises      []Exercise `json:"exercises,omitempty"`      // 训练项目列表
	TotalWeight    *float64   `json:"totalWeight,omitempty"`    // 总重量(kg)
	TotalSets      *int       `json:"totalSets,omitempty"`      // 总组数
	CaloriesBurned *int       `json:"caloriesBurned,omitempty"` // 消耗卡路里
	Notes          *string    `json:"notes,omitempty"`          // 训练备注
	Mood           *string    `json:"mood,omitempty"`           // 训练状态
	PlanID         *int       `json:"planId,omitempty"`         // 关联计划ID
}

// UpdateTrainingRecordRequest 更新训练记录请求
type UpdateTrainingRecordRequest struct {
	Title          *string    `json:"title,omitempty"`          // 标题
	StartTime      *string    `json:"startTime,omitempty"`      // 开始时间 YYYY-MM-DD HH:mm:ss
	EndTime        *string    `json:"endTime,omitempty"`        // 结束时间 YYYY-MM-DD HH:mm:ss
	Duration       *int       `json:"duration,omitempty"`       // 总时长(分钟)
	Exercises      []Exercise `json:"exercises,omitempty"`      // 训练项目列表
	TotalWeight    *float64   `json:"totalWeight,omitempty"`    // 总重量(kg)
	TotalSets      *int       `json:"totalSets,omitempty"`      // 总组数
	CaloriesBurned *int       `json:"caloriesBurned,omitempty"` // 消耗卡路里
	Notes          *string    `json:"notes,omitempty"`          // 训练备注
	Mood           *string    `json:"mood,omitempty"`           // 训练状态
	PlanID         *int       `json:"planId,omitempty"`         // 关联计划ID
}

// TrainingRecordUsecase 训练记录用例接口
type TrainingRecordUsecase interface {
	Create(c context.Context, userID string, request *CreateTrainingRecordRequest) (map[string]interface{}, error)
	GetByID(c context.Context, userID, recordID string) (TrainingRecord, error)
	GetList(c context.Context, userID string, page, pageSize int, startDate, endDate string, planID int) ([]TrainingRecord, int64, error)
	Update(c context.Context, userID, recordID string, request *UpdateTrainingRecordRequest) error
	Delete(c context.Context, userID, recordID string) error
}
