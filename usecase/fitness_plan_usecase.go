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
	template, err := fu.planTemplateRepository.GetByID(ctx, request.TemplateID)
	if err != nil {
		return nil, errors.New("template not found")
	}

	// Convert userID string to ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Convert templateID string to ObjectID
	templateObjectID, err := primitive.ObjectIDFromHex(request.TemplateID)
	if err != nil {
		return nil, errors.New("invalid template ID")
	}

	// Parse start date and calculate end date
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format")
	}

	// Use override duration if provided
	durationWeeks := template.DurationWeeks
	if request.DurationWeeksOverride != nil && *request.DurationWeeksOverride > 0 {
		durationWeeks = *request.DurationWeeksOverride
	}

	// Calculate end date based on duration
	endDate := startDate.AddDate(0, 0, durationWeeks*7-1)

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

	// Handle training days override
	var trainingDaysOverride []domain.TrainingDay
	if request.TrainingDaysOverride != nil && len(request.TrainingDaysOverride) > 0 {
		trainingDaysOverride = request.TrainingDaysOverride
	}

	completedDays := []int{}
	skippedDays := []int{}

	now := time.Now()
	plan := &domain.FitnessPlan{
		ID:                    primitive.NewObjectID(),
		UserID:                userObjectID,
		TemplateID:            &templateObjectID,
		Name:                  planName,
		Description:           template.Description,
		Goal:                  template.Goal,
		DurationWeeks:         durationWeeks,
		DurationWeeksOverride: request.DurationWeeksOverride,
		TrainingDaysPerWeek:   template.TrainingDaysPerWeek,
		TrainingDays:          trainingDays,
		TrainingDaysOverride:  trainingDaysOverride,
		StartDate:             request.StartDate,
		EndDate:               endDate.Format("2006-01-02"),
		Status:                "进行中",
		CurrentWeek:           1,
		CurrentDay:            1,
		CompletedDays:         completedDays,
		SkippedDays:           skippedDays,
		TotalCompletedDays:    0,
		CompletionRate:        0,
		TotalWeight:           0,
		TotalDuration:         0,
		TotalCalories:         0,
		CreatedAt:             primitive.NewDateTimeFromTime(now),
		UpdatedAt:             primitive.NewDateTimeFromTime(now),
	}

	err = fu.fitnessPlanRepository.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        plan.ID.Hex(),
		"name":      plan.Name,
		"startDate": plan.StartDate,
		"endDate":   plan.EndDate,
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
	skippedDays := []int{}

	now := time.Now()
	plan := &domain.FitnessPlan{
		ID:                  primitive.NewObjectID(),
		UserID:              userObjectID,
		TemplateID:          nil, // Custom plan has no template
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
		SkippedDays:         skippedDays,
		TotalCompletedDays:  0,
		CompletionRate:      0,
		TotalWeight:         0,
		TotalDuration:       0,
		TotalCalories:       0,
		CreatedAt:           primitive.NewDateTimeFromTime(now),
		UpdatedAt:           primitive.NewDateTimeFromTime(now),
	}

	err = fu.fitnessPlanRepository.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":        plan.ID.Hex(),
		"name":      plan.Name,
		"startDate": plan.StartDate,
		"endDate":   plan.EndDate,
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

func (fu *fitnessPlanUsecase) CompleteDay(c context.Context, userID, planID string, dayNumber int, recordID string) (map[string]interface{}, error) {
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

func (fu *fitnessPlanUsecase) UncompleteDay(c context.Context, userID, planID string, dayNumber int) (map[string]interface{}, error) {
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

	// Check if day is actually completed
	isCompleted := false
	for _, completedDay := range plan.CompletedDays {
		if completedDay == dayNumber {
			isCompleted = true
			break
		}
	}

	if !isCompleted {
		return nil, errors.New("day is not completed")
	}

	// Uncomplete the day
	err = fu.fitnessPlanRepository.UncompletePlanDay(ctx, planID, dayNumber)
	if err != nil {
		return nil, err
	}

	// Get updated plan for response
	updatedPlan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"completionRate":     updatedPlan.CompletionRate,
		"totalCompletedDays": updatedPlan.TotalCompletedDays,
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

func (fu *fitnessPlanUsecase) GetProgress(c context.Context, userID, planID string) (domain.PlanProgress, error) {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()

	// Get existing plan and validate ownership
	plan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return domain.PlanProgress{}, err
	}

	if plan.UserID.Hex() != userID {
		return domain.PlanProgress{}, errors.New("unauthorized access to fitness plan")
	}

	// Calculate total training days (non-rest days)
	totalDays := plan.DurationWeeks * plan.TrainingDaysPerWeek
	completedDays := len(plan.CompletedDays)
	skippedDays := len(plan.SkippedDays)

	// Calculate completion rate
	effectiveTotalDays := totalDays - skippedDays
	completionRate := 0
	if effectiveTotalDays > 0 {
		completionRate = (completedDays * 100) / effectiveTotalDays
	}

	// Calculate next training date
	nextTrainingDate := ""
	if plan.Status == "进行中" {
		startDate, err := time.Parse("2006-01-02", plan.StartDate)
		if err == nil {
			// Find next uncompleted day
			nextDayNum := completedDays + skippedDays + 1
			// Calculate the date for the next day
			nextDate := startDate.AddDate(0, 0, nextDayNum-1)
			nextTrainingDate = nextDate.Format("2006-01-02")
		}
	}

	// Calculate current week and day
	currentWeek := 1
	currentDay := 1
	if plan.Status == "进行中" {
		startDate, err := time.Parse("2006-01-02", plan.StartDate)
		if err == nil {
			daysSinceStart := int(time.Since(startDate).Hours() / 24)
			if daysSinceStart >= 0 {
				currentWeek = (daysSinceStart / 7) + 1
				currentDay = (daysSinceStart % 7) + 1
			}
		}
	}

	progress := domain.PlanProgress{
		PlanID:           planID,
		TotalDays:        totalDays,
		CompletedDays:    completedDays,
		SkippedDays:      skippedDays,
		CompletionRate:   completionRate,
		CurrentWeek:      currentWeek,
		CurrentDay:       currentDay,
		NextTrainingDate: nextTrainingDate,
		TotalDuration:    plan.TotalDuration,
		TotalWeight:      plan.TotalWeight,
		TotalCalories:    plan.TotalCalories,
	}

	return progress, nil
}

func (fu *fitnessPlanUsecase) SkipDay(c context.Context, userID, planID string, dayNumber int, reason string) (map[string]interface{}, error) {
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
			return nil, errors.New("day already completed, cannot skip")
		}
	}

	// Check if day is already skipped
	for _, skippedDay := range plan.SkippedDays {
		if skippedDay == dayNumber {
			return nil, errors.New("day already skipped")
		}
	}

	// Skip the day
	err = fu.fitnessPlanRepository.SkipPlanDay(ctx, planID, dayNumber)
	if err != nil {
		return nil, err
	}

	// Get updated plan
	updatedPlan, err := fu.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"skippedDays":    len(updatedPlan.SkippedDays),
		"completionRate": updatedPlan.CompletionRate,
	}, nil
}

func (fu *fitnessPlanUsecase) AdjustDay(c context.Context, userID, planID string, dayNumber int, exercises []domain.Exercise, notes string) error {
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

	// Validate day number
	totalDays := plan.DurationWeeks * 7
	if dayNumber < 1 || dayNumber > totalDays {
		return errors.New("invalid day number")
	}

	return fu.fitnessPlanRepository.UpdateTrainingDay(ctx, planID, dayNumber, exercises, notes)
}
