package repository

import (
	"context"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
	"github.com/zhengshui/flow-link-server/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type planTemplateRepository struct {
	database   mongo.Database
	collection string
}

func NewPlanTemplateRepository(db mongo.Database, collection string) domain.PlanTemplateRepository {
	return &planTemplateRepository{
		database:   db,
		collection: collection,
	}
}

func (pt *planTemplateRepository) GetByID(c context.Context, id string) (domain.PlanTemplate, error) {
	collection := pt.database.Collection(pt.collection)
	var template domain.PlanTemplate

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return template, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&template)
	return template, err
}

func (pt *planTemplateRepository) GetList(c context.Context, goal, level string, page, pageSize int) ([]domain.PlanTemplate, int64, error) {
	collection := pt.database.Collection(pt.collection)

	// 构建查询条件
	filter := bson.M{}
	if goal != "" {
		filter["goal"] = goal
	}
	if level != "" {
		filter["level"] = level
	}

	// 计算总数
	total, err := collection.CountDocuments(c, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(c, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var templates []domain.PlanTemplate
	err = cursor.All(c, &templates)
	if templates == nil {
		return []domain.PlanTemplate{}, total, err
	}

	return templates, total, err
}

func (pt *planTemplateRepository) Create(c context.Context, template *domain.PlanTemplate) error {
	collection := pt.database.Collection(pt.collection)
	template.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := collection.InsertOne(c, template)
	return err
}
