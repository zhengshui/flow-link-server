package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionFeedback = "feedbacks"
)

// Feedback 用户反馈模型
type Feedback struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Content     string             `bson:"content" json:"content"`
	Type        string             `bson:"type" json:"type"`               // 反馈类型：建议/问题/其他
	ContactInfo string             `bson:"contactInfo" json:"contactInfo"` // 联系方式（可选）
	Status      string             `bson:"status" json:"status"`           // 处理状态：待处理/处理中/已处理
	CreatedAt   primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

// FeedbackRequest 创建反馈请求
type FeedbackRequest struct {
	Content     string `json:"content" binding:"required,min=1,max=1000"` // 反馈内容（必填，10-1000字符）
	Type        string `json:"type"`                                      // 反馈类型：建议/问题/其他（可选，默认：建议）
	ContactInfo string `json:"contactInfo"`                               // 联系方式（可选）
}

// FeedbackResponse 创建反馈响应
type FeedbackResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

// FeedbackRepository 反馈仓库接口
type FeedbackRepository interface {
	Create(c context.Context, feedback *Feedback) error
}

// FeedbackUsecase 反馈用例接口
type FeedbackUsecase interface {
	Create(c context.Context, feedback *Feedback) error
}
