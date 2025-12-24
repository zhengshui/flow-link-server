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

type trainingRecordRepository struct {
	database   mongo.Database
	collection string
}

func NewTrainingRecordRepository(db mongo.Database, collection string) domain.TrainingRecordRepository {
	return &trainingRecordRepository{
		database:   db,
		collection: collection,
	}
}

func (tr *trainingRecordRepository) Create(c context.Context, record *domain.TrainingRecord) error {
	collection := tr.database.Collection(tr.collection)
	record.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	record.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := collection.InsertOne(c, record)
	return err
}

func (tr *trainingRecordRepository) GetByID(c context.Context, id string) (domain.TrainingRecord, error) {
	collection := tr.database.Collection(tr.collection)
	var record domain.TrainingRecord

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return record, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&record)
	return record, err
}

func (tr *trainingRecordRepository) GetByUserID(c context.Context, userID string, page, pageSize int, startDate, endDate string, planID string) ([]domain.TrainingRecord, int64, error) {
	collection := tr.database.Collection(tr.collection)

	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}

	// 构建查询条件
	filter := bson.M{"userId": userIDHex}

	if startDate != "" && endDate != "" {
		// 如果日期格式是 YYYY-MM-DD，补充时间部分以确保查询正确
		if len(startDate) == 10 {
			startDate = startDate + " 00:00:00"
		}
		if len(endDate) == 10 {
			endDate = endDate + " 23:59:59"
		}

		filter["startTime"] = bson.M{
			"$gte": startDate,
			"$lte": endDate,
		}
	}

	if planID != "" {
		filter["planId"] = planID
	}

	// 计算总数
	total, err := collection.CountDocuments(c, filter)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSort(bson.D{{Key: "startTime", Value: -1}, {Key: "createdAt", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := collection.Find(c, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var records []domain.TrainingRecord
	err = cursor.All(c, &records)
	if records == nil {
		return []domain.TrainingRecord{}, total, err
	}

	return records, total, err
}

func (tr *trainingRecordRepository) Update(c context.Context, id string, record *domain.TrainingRecord) error {
	collection := tr.database.Collection(tr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	record.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": bson.M{
			"title":          record.Title,
			"startTime":      record.StartTime,
			"endTime":        record.EndTime,
			"duration":       record.Duration,
			"exercises":      record.Exercises,
			"totalWeight":    record.TotalWeight,
			"totalSets":      record.TotalSets,
			"caloriesBurned": record.CaloriesBurned,
			"notes":          record.Notes,
			"mood":           record.Mood,
			"planId":         record.PlanID,
			"updatedAt":      record.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (tr *trainingRecordRepository) Delete(c context.Context, id string) error {
	collection := tr.database.Collection(tr.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(c, bson.M{"_id": idHex})
	return err
}
