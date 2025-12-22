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

func (pt *planTemplateRepository) GetList(c context.Context, goal, level, splitType, equipment string, durationWeeksMin, durationWeeksMax, page, pageSize int) ([]domain.PlanTemplate, int64, error) {
	collection := pt.database.Collection(pt.collection)

	// 构建查询条件 - 只查询官方模板
	filter := bson.M{"isOfficial": true}
	if goal != "" {
		filter["goal"] = goal
	}
	if level != "" {
		filter["level"] = level
	}
	if splitType != "" {
		filter["splitType"] = splitType
	}
	if equipment != "" {
		filter["equipment"] = equipment
	}
	if durationWeeksMin > 0 {
		filter["durationWeeks"] = bson.M{"$gte": durationWeeksMin}
	}
	if durationWeeksMax > 0 {
		if _, exists := filter["durationWeeks"]; exists {
			filter["durationWeeks"].(bson.M)["$lte"] = durationWeeksMax
		} else {
			filter["durationWeeks"] = bson.M{"$lte": durationWeeksMax}
		}
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

func (pt *planTemplateRepository) GetUserTemplates(c context.Context, userID string, page, pageSize int) ([]domain.PlanTemplate, int64, error) {
	collection := pt.database.Collection(pt.collection)

	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}

	// 查询用户的个人模板
	filter := bson.M{
		"userId":     userIDHex,
		"isOfficial": false,
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
	now := time.Now()
	template.CreatedAt = primitive.NewDateTimeFromTime(now)
	template.UpdatedAt = primitive.NewDateTimeFromTime(now)
	_, err := collection.InsertOne(c, template)
	return err
}

func (pt *planTemplateRepository) Update(c context.Context, id string, template *domain.PlanTemplate) error {
	collection := pt.database.Collection(pt.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	template.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": bson.M{
			"name":                 template.Name,
			"description":          template.Description,
			"goal":                 template.Goal,
			"splitType":            template.SplitType,
			"level":                template.Level,
			"equipment":            template.Equipment,
			"durationWeeks":        template.DurationWeeks,
			"trainingDaysPerWeek":  template.TrainingDaysPerWeek,
			"trainingDays":         template.TrainingDays,
			"tags":                 template.Tags,
			"imageUrl":             template.ImageUrl,
			"recommendedIntensity": template.RecommendedIntensity,
			"updatedAt":            template.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (pt *planTemplateRepository) Delete(c context.Context, id string) error {
	collection := pt.database.Collection(pt.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(c, bson.M{"_id": idHex})
	return err
}
