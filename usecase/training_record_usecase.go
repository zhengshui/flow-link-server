package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trainingRecordUsecase struct {
	trainingRecordRepository domain.TrainingRecordRepository
	contextTimeout           time.Duration
}

func NewTrainingRecordUsecase(trainingRecordRepository domain.TrainingRecordRepository, timeout time.Duration) domain.TrainingRecordUsecase {
	return &trainingRecordUsecase{
		trainingRecordRepository: trainingRecordRepository,
		contextTimeout:           timeout,
	}
}

func (tu *trainingRecordUsecase) Create(c context.Context, userID string, request *domain.CreateTrainingRecordRequest) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Initialize exercises array to avoid null in JSON
	exercises := request.Exercises
	if exercises == nil {
		exercises = []domain.Exercise{}
	}

	now := time.Now()
	record := &domain.TrainingRecord{
		ID:             primitive.NewObjectID(),
		UserID:         userObjectID,
		Title:          request.Title,
		StartTime:      request.StartTime,
		EndTime:        request.EndTime,
		Duration:       request.Duration,
		Exercises:      exercises,
		TotalWeight:    request.TotalWeight,
		TotalSets:      request.TotalSets,
		CaloriesBurned: request.CaloriesBurned,
		Notes:          request.Notes,
		Mood:           request.Mood,
		PlanID:         request.PlanID,
		CreatedAt:      primitive.NewDateTimeFromTime(now),
		UpdatedAt:      primitive.NewDateTimeFromTime(now),
	}

	err = tu.trainingRecordRepository.Create(ctx, record)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        record.ID.Hex(),
		"createdAt": record.CreatedAt,
	}, nil
}

func (tu *trainingRecordUsecase) GetByID(c context.Context, userID, recordID string) (domain.TrainingRecord, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	record, err := tu.trainingRecordRepository.GetByID(ctx, recordID)
	if err != nil {
		return domain.TrainingRecord{}, err
	}

	// Validate ownership
	if record.UserID.Hex() != userID {
		return domain.TrainingRecord{}, errors.New("unauthorized access to training record")
	}

	return record, nil
}

func (tu *trainingRecordUsecase) GetList(c context.Context, userID string, page, pageSize int, startDate, endDate string, planID int) ([]domain.TrainingRecord, int64, error) {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	records, total, err := tu.trainingRecordRepository.GetByUserID(ctx, userID, page, pageSize, startDate, endDate, planID)
	if err != nil {
		return nil, 0, err
	}

	// Initialize empty array to avoid null in JSON
	if records == nil {
		records = []domain.TrainingRecord{}
	}

	return records, total, nil
}

func (tu *trainingRecordUsecase) Update(c context.Context, userID, recordID string, request *domain.UpdateTrainingRecordRequest) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// Get existing record and validate ownership
	record, err := tu.trainingRecordRepository.GetByID(ctx, recordID)
	if err != nil {
		return err
	}

	if record.UserID.Hex() != userID {
		return errors.New("unauthorized access to training record")
	}

	// Update fields if provided
	if request.Title != "" {
		record.Title = request.Title
	}
	if request.StartTime != "" {
		record.StartTime = request.StartTime
	}
	if request.EndTime != "" {
		record.EndTime = request.EndTime
	}
	if request.Duration > 0 {
		record.Duration = request.Duration
	}
	if request.Exercises != nil {
		record.Exercises = request.Exercises
	}
	if request.TotalWeight >= 0 {
		record.TotalWeight = request.TotalWeight
	}
	if request.TotalSets >= 0 {
		record.TotalSets = request.TotalSets
	}
	if request.CaloriesBurned >= 0 {
		record.CaloriesBurned = request.CaloriesBurned
	}
	if request.Notes != "" {
		record.Notes = request.Notes
	}
	if request.Mood != "" {
		record.Mood = request.Mood
	}
	if request.PlanID >= 0 {
		record.PlanID = request.PlanID
	}

	record.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return tu.trainingRecordRepository.Update(ctx, recordID, &record)
}

func (tu *trainingRecordUsecase) Delete(c context.Context, userID, recordID string) error {
	ctx, cancel := context.WithTimeout(c, tu.contextTimeout)
	defer cancel()

	// Get existing record and validate ownership
	record, err := tu.trainingRecordRepository.GetByID(ctx, recordID)
	if err != nil {
		return err
	}

	if record.UserID.Hex() != userID {
		return errors.New("unauthorized access to training record")
	}

	return tu.trainingRecordRepository.Delete(ctx, recordID)
}
