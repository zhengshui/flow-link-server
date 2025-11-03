package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type fitnessPlanUsecase struct {
	fitnessPlanRepository  domain.FitnessPlanRepository
	planTemplateRepository domain.PlanTemplateRepository
	contextTimeout         time.Duration
}

func NewFitnessPlanUsecase(fitnessPlanRepository domain.FitnessPlanRepository, planTemplateRepository domain.PlanTemplateRepository, timeout time.Duration) domain.FitnessPlanUsecase {
	return &fitnessPlanUsecase{
		fitnessPlanRepository:  fitnessPlanRepository,
		planTemplateRepository: planTemplateRepository,
		contextTimeout:         timeout,
	}
}

func (fu *fitnessPlanUsecase) CreateFromTemplate(c context.Context, userID string, request *domain.CreatePlanFromTemplateRequest) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	// Get template
	template, err := fu.planTemplateRepository.GetByID(ctx, string(rune(request.TemplateID)))
	if err != nil {
		return nil, errors.New("template not found")
	}

	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Parse start date and calculate end date
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	// Calculate end date based on duration
	endDate := startDate.AddDate(0, 0, template.DurationWeeks*7-1)

	// Use custom name if provided, otherwise use template name
	planName := request.Name
	if planName == "" {
		planName = template.Name
	}

	// Initialize arrays to avoid null in JSON
	trainingDays := template.TrainingDays
	if trainingDays == nil {
		trainingDays = []domain.TrainingDay{}
	}

	completedDays := []int{}

	now := time.Now()
	plan := &domain.FitnessPlan{
		ID:                  primitive.NewObjectID(),
		UserID:              userObjectID,
		TemplateID:          request.TemplateID,
		Name:                planName,
		Description:         template.Description,
		Goal:                template.Goal,
		DurationWeeks:       template.DurationWeeks,
		TrainingDaysPerWeek: template.TrainingDaysPerWeek,
		TrainingDays:        trainingDays,
		StartDate:           request.StartDate,
		EndDate:             endDate.Format("2006-01-02"),
		Status:              "进行中",
		CurrentWeek:         1,
		CurrentDay:          1,
		CompletedDays:       completedDays,
		TotalCompletedDays:  0,
		CompletionRate:      0,
		CreatedAt:           primitive.NewDateTimeFromTime(now),
		UpdatedAt:           primitive.NewDateTimeFromTime(now),
	}

	err = fu.fitnessPlanRepository.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        plan.ID.Hex(),
		"createdAt": plan.CreatedAt,
	}, nil
}

func (fu *fitnessPlanUsecase) CreateCustom(c context.Context, userID string, request *domain.CreateCustomPlanRequest) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Parse start date and calculate end date
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	// Calculate end date based on duration
	endDate := startDate.AddDate(0, 0, request.DurationWeeks*7-1)

	// Initialize arrays to avoid null in JSON
	trainingDays := request.TrainingDays
	if trainingDays == nil {
		trainingDays = []domain.TrainingDay{}
	}

	completedDays := []int{}

	now := time.Now()
	plan := &domain.FitnessPlan{
		ID:                  primitive.NewObjectID(),
		UserID:              userObjectID,
		TemplateID:          0, // Custom plan has no template
		Name:                request.Name,
		Description:         request.Description,
		Goal:                request.Goal,
		DurationWeeks:       request.DurationWeeks,
		TrainingDaysPerWeek: request.TrainingDaysPerWeek,
		TrainingDays:        trainingDays,
		StartDate:           request.StartDate,
		EndDate:             endDate.Format("2006-01-02"),
		Status:              "进行中",
		CurrentWeek:         1,
		CurrentDay:          1,
		CompletedDays:       completedDays,
		TotalCompletedDays:  0,
		CompletionRate:      0,
		CreatedAt:           primitive.NewDateTimeFromTime(now),
		UpdatedAt:           primitive.NewDateTimeFromTime(now),
	}

	err = fu.fitnessPlanRepository.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        plan.ID.Hex(),
		"createdAt": plan.CreatedAt,
	}, nil
}

func (fu *fitnessPlanUsecase) GetByID(c context.Context, userID, planID string) (domain.FitnessPlan, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	plan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return domain.FitnessPlan{}, err
	}

	// Validate ownership
	if plan.UserID.Hex() != userID {
		return domain.FitnessPlan{}, errors.New("unauthorized access to fitness plan")
	}

	return plan, nil
}

func (fu *fitnessPlanUsecase) GetList(c context.Context, userID string, status string, page, pageSize int) ([]domain.FitnessPlan, int64, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	plans, total, err := fu.fitnessPlanRepository.GetByUserID(ctx, userID, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Initialize empty array to avoid null in JSON
	if plans == nil {
		plans = []domain.FitnessPlan{}
	}

	return plans, total, nil
}

func (fu *fitnessPlanUsecase) UpdateStatus(c context.Context, userID, planID string, status string) error {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	// Get existing plan and validate ownership
	plan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	if plan.UserID.Hex() != userID {
		return errors.New("unauthorized access to fitness plan")
	}

	return fu.fitnessPlanRepository.UpdateStatus(ctx, planID, status)
}

func (fu *fitnessPlanUsecase) CompleteDay(c context.Context, userID, planID string, dayNumber int) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	// Get existing plan and validate ownership
	plan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	if plan.UserID.Hex() != userID {
		return nil, errors.New("unauthorized access to fitness plan")
	}

	// Check if day is already completed
	for _, completedDay := range plan.CompletedDays {
		if completedDay == dayNumber {
			return nil, errors.New("day already completed")
		}
	}

	// Add day to completed days
	plan.CompletedDays = append(plan.CompletedDays, dayNumber)
	plan.TotalCompletedDays++

	// Calculate total training days (non-rest days)
	totalTrainingDays := 0
	for _, day := range plan.TrainingDays {
		if !day.IsRestDay {
			totalTrainingDays++
		}
	}

	// Calculate completion rate
	if totalTrainingDays > 0 {
		plan.CompletionRate = (plan.TotalCompletedDays * 100) / (totalTrainingDays * plan.DurationWeeks)
	}

	plan.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Update the plan
	err = fu.fitnessPlanRepository.Update(ctx, planID, &plan)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"completionRate":     plan.CompletionRate,
		"totalCompletedDays": plan.TotalCompletedDays,
	}, nil
}

func (fu *fitnessPlanUsecase) Delete(c context.Context, userID, planID string) error {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	// Get existing plan and validate ownership
	plan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	if plan.UserID.Hex() != userID {
		return errors.New("unauthorized access to fitness plan")
	}

	return fu.fitnessPlanRepository.Delete(ctx, planID)
}
