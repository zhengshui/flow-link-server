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
	ID          int     `bson:"id" json:"id"`
	Name        string  `bson:"name" json:"name"`                   // 项目名称
	Sets        int     `bson:"sets" json:"sets"`                   // 组数
	Reps        int     `bson:"reps" json:"reps"`                   // 次数
	Weight      float64 `bson:"weight" json:"weight"`               // 重量(kg)
	RestTime    int     `bson:"restTime" json:"restTime"`           // 休息时间(秒)
	MuscleGroup string  `bson:"muscleGroup" json:"muscleGroup"`     // 目标肌群
	Notes       string  `bson:"notes" json:"notes,omitempty"`       // 备注
	Duration    int     `bson:"duration" json:"duration,omitempty"` // 训练时长(分钟)
}

// TrainingRecord 训练记录
type TrainingRecord struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Title          string             `bson:"title" json:"title"`
	StartTime      string             `bson:"startTime" json:"startTime"`           // 开始时间 YYYY-MM-DD HH:mm:ss
	EndTime        string             `bson:"endTime" json:"endTime"`               // 结束时间 YYYY-MM-DD HH:mm:ss
	Duration       int                `bson:"duration" json:"duration"`             // 总时长(分钟)
	Exercises      []Exercise         `bson:"exercises" json:"exercises"`           // 训练项目列表
	TotalWeight    float64            `bson:"totalWeight" json:"totalWeight"`       // 总重量(kg)
	TotalSets      int                `bson:"totalSets" json:"totalSets"`           // 总组数
	CaloriesBurned int                `bson:"caloriesBurned" json:"caloriesBurned"` // 消耗卡路里
	Notes          string             `bson:"notes" json:"notes,omitempty"`         // 训练备注
	Mood           string             `bson:"mood" json:"mood,omitempty"`           // 训练状态(优秀/良好/一般/疲劳)
	PlanID         int                `bson:"planId" json:"planId"`                 // 关联计划ID(0表示无计划)
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
	Title          string     `json:"title" binding:"required"`
	StartTime      string     `json:"startTime" binding:"required"` // YYYY-MM-DD HH:mm:ss
	EndTime        string     `json:"endTime" binding:"required"`   // YYYY-MM-DD HH:mm:ss
	Duration       int        `json:"duration" binding:"required"`
	Exercises      []Exercise `json:"exercises" binding:"required"`
	TotalWeight    float64    `json:"totalWeight"`
	TotalSets      int        `json:"totalSets"`
	CaloriesBurned int        `json:"caloriesBurned"`
	Notes          string     `json:"notes"`
	Mood           string     `json:"mood"`
	PlanID         int        `json:"planId"`
}

// UpdateTrainingRecordRequest 更新训练记录请求
type UpdateTrainingRecordRequest struct {
	Title          string     `json:"title"`
	StartTime      string     `json:"startTime"` // YYYY-MM-DD HH:mm:ss
	EndTime        string     `json:"endTime"`   // YYYY-MM-DD HH:mm:ss
	Duration       int        `json:"duration"`
	Exercises      []Exercise `json:"exercises"`
	TotalWeight    float64    `json:"totalWeight"`
	TotalSets      int        `json:"totalSets"`
	CaloriesBurned int        `json:"caloriesBurned"`
	Notes          string     `json:"notes"`
	Mood           string     `json:"mood"`
	PlanID         int        `json:"planId"`
}

// TrainingRecordUsecase 训练记录用例接口
type TrainingRecordUsecase interface {
	Create(c context.Context, userID string, request *CreateTrainingRecordRequest) (map[string]interface{}, error)
	GetByID(c context.Context, userID, recordID string) (TrainingRecord, error)
	GetList(c context.Context, userID string, page, pageSize int, startDate, endDate string, planID int) ([]TrainingRecord, int64, error)
	Update(c context.Context, userID, recordID string, request *UpdateTrainingRecordRequest) error
	Delete(c context.Context, userID, recordID string) error
}
