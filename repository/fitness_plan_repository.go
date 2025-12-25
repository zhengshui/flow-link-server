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

type fitnessPlanRepository struct {
	database   mongo.Database
	collection string
}

func NewFitnessPlanRepository(db mongo.Database, collection string) domain.FitnessPlanRepository {
	return &fitnessPlanRepository{
		database:   db,
		collection: collection,
	}
}

func (fp *fitnessPlanRepository) Create(c context.Context, plan *domain.FitnessPlan) error {
	collection := fp.database.Collection(fp.collection)
	plan.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	plan.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := collection.InsertOne(c, plan)
	return err
}

func (fp *fitnessPlanRepository) GetByID(c context.Context, id string) (domain.FitnessPlan, error) {
	collection := fp.database.Collection(fp.collection)
	var plan domain.FitnessPlan

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return plan, err
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&plan)
	return plan, err
}

func (fp *fitnessPlanRepository) GetByUserID(c context.Context, userID string, status string, page, pageSize int) ([]domain.FitnessPlan, int64, error) {
	collection := fp.database.Collection(fp.collection)

	userIDHex, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}

	// 构建查询条件
	filter := bson.M{"userId": userIDHex}
	if status != "" {
		filter["status"] = status
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

	var plans []domain.FitnessPlan
	err = cursor.All(c, &plans)
	if plans == nil {
		return []domain.FitnessPlan{}, total, err
	}

	return plans, total, err
}

func (fp *fitnessPlanRepository) Update(c context.Context, id string, plan *domain.FitnessPlan) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	plan.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": bson.M{
			"name":                plan.Name,
			"description":         plan.Description,
			"goal":                plan.Goal,
			"durationWeeks":       plan.DurationWeeks,
			"trainingDaysPerWeek": plan.TrainingDaysPerWeek,
			"trainingDays":        plan.TrainingDays,
			"startDate":           plan.StartDate,
			"endDate":             plan.EndDate,
			"status":              plan.Status,
			"currentWeek":         plan.CurrentWeek,
			"currentDay":          plan.CurrentDay,
			"completedDays":       plan.CompletedDays,
			"totalCompletedDays":  plan.TotalCompletedDays,
			"completionRate":      plan.CompletionRate,
			"updatedAt":           plan.UpdatedAt,
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (fp *fitnessPlanRepository) UpdateStatus(c context.Context, id string, status string) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (fp *fitnessPlanRepository) Delete(c context.Context, id string) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(c, bson.M{"_id": idHex})
	return err
}

func (fp *fitnessPlanRepository) CompletePlanDay(c context.Context, id string, dayNumber int) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 添加dayNumber到completedDays数组，如果不存在的话
	update := bson.M{
		"$addToSet": bson.M{
			"completedDays": dayNumber,
		},
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	if err != nil {
		return err
	}

	// 获取更新后的计划，重新计算totalCompletedDays和completionRate
	var plan domain.FitnessPlan
	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&plan)
	if err != nil {
		return err
	}

	totalDays := plan.DurationWeeks * plan.TrainingDaysPerWeek
	totalCompletedDays := len(plan.CompletedDays)
	completionRate := 0
	if totalDays > 0 {
		completionRate = (totalCompletedDays * 100) / totalDays
	}

	// 更新统计数据
	update = bson.M{
		"$set": bson.M{
			"totalCompletedDays": totalCompletedDays,
			"completionRate":     completionRate,
			"updatedAt":          primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (fp *fitnessPlanRepository) UncompletePlanDay(c context.Context, id string, dayNumber int) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Remove dayNumber from completedDays array
	update := bson.M{
		"$pull": bson.M{
			"completedDays": dayNumber,
		},
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	if err != nil {
		return err
	}

	// Recalculate stats
	var plan domain.FitnessPlan
	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&plan)
	if err != nil {
		return err
	}

	totalDays := plan.DurationWeeks * plan.TrainingDaysPerWeek
	totalCompletedDays := len(plan.CompletedDays)
	completionRate := 0
	if totalDays > 0 {
		completionRate = (totalCompletedDays * 100) / totalDays
	}

	// Update stats
	update = bson.M{
		"$set": bson.M{
			"totalCompletedDays": totalCompletedDays,
			"completionRate":     completionRate,
			"updatedAt":          primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (fp *fitnessPlanRepository) SkipPlanDay(c context.Context, id string, dayNumber int) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 添加dayNumber到skippedDays数组，如果不存在的话
	update := bson.M{
		"$addToSet": bson.M{
			"skippedDays": dayNumber,
		},
		"$set": bson.M{
			"updatedAt": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	if err != nil {
		return err
	}

	// 获取更新后的计划，重新计算completionRate
	var plan domain.FitnessPlan
	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&plan)
	if err != nil {
		return err
	}

	totalDays := plan.DurationWeeks * plan.TrainingDaysPerWeek
	skippedDays := len(plan.SkippedDays)
	completedDays := len(plan.CompletedDays)
	// 完成率 = (已完成 / (总天数 - 跳过天数)) * 100
	effectiveTotalDays := totalDays - skippedDays
	completionRate := 0
	if effectiveTotalDays > 0 {
		completionRate = (completedDays * 100) / effectiveTotalDays
	}

	// 更新统计数据
	update = bson.M{
		"$set": bson.M{
			"completionRate": completionRate,
			"updatedAt":      primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}

func (fp *fitnessPlanRepository) UpdateTrainingDay(c context.Context, id string, dayNumber int, exercises []domain.Exercise, notes string) error {
	collection := fp.database.Collection(fp.collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 获取当前计划
	var plan domain.FitnessPlan
	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&plan)
	if err != nil {
		return err
	}

	// 更新 trainingDaysOverride
	trainingDaysOverride := plan.TrainingDaysOverride
	if trainingDaysOverride == nil {
		trainingDaysOverride = []domain.TrainingDay{}
	}

	// 查找是否已存在该天的覆盖
	found := false
	for i, day := range trainingDaysOverride {
		if day.DayNumber == dayNumber {
			trainingDaysOverride[i].Exercises = exercises
			if notes != "" {
				trainingDaysOverride[i].Notes = notes
			}
			found = true
			break
		}
	}

	// 如果不存在，则添加新的覆盖
	if !found {
		// 从原始训练日中获取基础信息
		var baseDayName string
		var isRestDay bool
		for _, day := range plan.TrainingDays {
			if day.DayNumber == dayNumber {
				baseDayName = day.DayName
				isRestDay = day.IsRestDay
				break
			}
		}
		newDay := domain.TrainingDay{
			DayNumber: dayNumber,
			DayName:   baseDayName,
			IsRestDay: isRestDay,
			Exercises: exercises,
			Notes:     notes,
		}
		trainingDaysOverride = append(trainingDaysOverride, newDay)
	}

	// 更新数据库
	update := bson.M{
		"$set": bson.M{
			"trainingDaysOverride": trainingDaysOverride,
			"updatedAt":            primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	_, err = collection.UpdateOne(c, bson.M{"_id": idHex}, update)
	return err
}
