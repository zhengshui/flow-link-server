package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
)

type statsUsecase struct {
	trainingRecordRepository domain.TrainingRecordRepository
	fitnessPlanRepository    domain.FitnessPlanRepository
	contextTimeout           time.Duration
}

func NewStatsUsecase(trainingRecordRepository domain.TrainingRecordRepository, fitnessPlanRepository domain.FitnessPlanRepository, timeout time.Duration) domain.StatsUsecase {
	return &statsUsecase{
		trainingRecordRepository: trainingRecordRepository,
		fitnessPlanRepository:    fitnessPlanRepository,
		contextTimeout:           timeout,
	}
}

func (su *statsUsecase) GetTrainingStats(c context.Context, userID string, period, startDate, endDate string) (domain.TrainingStats, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Calculate date range based on period if not provided
	if startDate == "" || endDate == "" {
		now := time.Now()
		switch period {
		case "week":
			startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
			endDate = now.Format("2006-01-02")
		case "month":
			startDate = now.AddDate(0, -1, 0).Format("2006-01-02")
			endDate = now.Format("2006-01-02")
		case "year":
			startDate = now.AddDate(-1, 0, 0).Format("2006-01-02")
			endDate = now.Format("2006-01-02")
		default:
			period = "week"
			startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
			endDate = now.Format("2006-01-02")
		}
	}

	// Get training records for the period
	records, _, err := su.trainingRecordRepository.GetByUserID(ctx, userID, 1, 1000, startDate, endDate, 0)
	if err != nil {
		return domain.TrainingStats{}, err
	}

	// Initialize stats
	stats := domain.TrainingStats{
		UserID:     userID,
		Period:     period,
		StartDate:  startDate,
		EndDate:    endDate,
		DailyStats: []domain.DailyStats{},
	}

	if len(records) == 0 {
		return stats, nil
	}

	// Calculate statistics
	totalTrainingCount := len(records)
	totalDuration := 0
	totalWeight := 0.0
	totalSets := 0
	totalCalories := 0

	// Maps for muscle group and exercise tracking
	muscleGroupCount := make(map[string]int)
	exerciseCount := make(map[string]int)
	dailyStatsMap := make(map[string]*domain.DailyStats)

	for _, record := range records {
		// 安全地累加指针类型的值
		if record.Duration != nil {
			totalDuration += *record.Duration
		}
		if record.TotalWeight != nil {
			totalWeight += *record.TotalWeight
		}
		if record.TotalSets != nil {
			totalSets += *record.TotalSets
		}
		if record.CaloriesBurned != nil {
			totalCalories += *record.CaloriesBurned
		}

		// Extract date from startTime (YYYY-MM-DD HH:mm:ss -> YYYY-MM-DD)
		recordDate := ""
		if record.StartTime != nil && len(*record.StartTime) >= 10 {
			recordDate = (*record.StartTime)[:10]
		}

		// 跳过没有有效日期的记录
		if recordDate == "" {
			continue
		}

		// Track daily stats
		if _, exists := dailyStatsMap[recordDate]; !exists {
			dailyStatsMap[recordDate] = &domain.DailyStats{
				Date: recordDate,
			}
		}
		dailyStats := dailyStatsMap[recordDate]
		dailyStats.TrainingCount++
		if record.Duration != nil {
			dailyStats.Duration += *record.Duration
		}
		if record.TotalWeight != nil {
			dailyStats.Weight += *record.TotalWeight
		}
		if record.TotalSets != nil {
			dailyStats.Sets += *record.TotalSets
		}
		if record.CaloriesBurned != nil {
			dailyStats.Calories += *record.CaloriesBurned
		}

		// Track muscle groups and exercises
		for _, exercise := range record.Exercises {
			if exercise.MuscleGroup != nil && *exercise.MuscleGroup != "" {
				muscleGroupCount[*exercise.MuscleGroup]++
			}
			if exercise.Name != "" {
				exerciseCount[exercise.Name]++
			}
		}
	}

	// Convert daily stats map to slice
	for _, ds := range dailyStatsMap {
		stats.DailyStats = append(stats.DailyStats, *ds)
	}

	// Find most trained muscle group
	maxMuscleCount := 0
	mostTrainedMuscle := ""
	for muscle, count := range muscleGroupCount {
		if count > maxMuscleCount {
			maxMuscleCount = count
			mostTrainedMuscle = muscle
		}
	}

	// Find favorite exercise
	maxExerciseCount := 0
	favoriteExercise := ""
	for exercise, count := range exerciseCount {
		if count > maxExerciseCount {
			maxExerciseCount = count
			favoriteExercise = exercise
		}
	}

	// Set calculated values
	stats.TotalTrainingCount = totalTrainingCount
	stats.TotalDuration = totalDuration
	stats.TotalWeight = totalWeight
	stats.TotalSets = totalSets
	stats.TotalCalories = totalCalories

	if totalTrainingCount > 0 {
		stats.AvgDuration = totalDuration / totalTrainingCount
		stats.AvgWeight = totalWeight / float64(totalTrainingCount)
	}

	stats.MostTrainedMuscle = mostTrainedMuscle
	stats.FavoriteExercise = favoriteExercise

	return stats, nil
}

func (su *statsUsecase) GetMuscleGroupStats(c context.Context, userID string, period string) ([]domain.MuscleGroupStats, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Calculate date range based on period
	now := time.Now()
	var startDate, endDate string

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	case "month":
		startDate = now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	case "year":
		startDate = now.AddDate(-1, 0, 0).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	default:
		startDate = now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	}

	// Get training records for the period
	records, _, err := su.trainingRecordRepository.GetByUserID(ctx, userID, 1, 1000, startDate, endDate, 0)
	if err != nil {
		return nil, err
	}

	// Calculate muscle group statistics
	muscleGroupData := make(map[string]*domain.MuscleGroupStats)
	totalCount := 0

	for _, record := range records {
		for _, exercise := range record.Exercises {
			if exercise.MuscleGroup != nil && *exercise.MuscleGroup != "" {
				muscleGroup := *exercise.MuscleGroup
				if _, exists := muscleGroupData[muscleGroup]; !exists {
					muscleGroupData[muscleGroup] = &domain.MuscleGroupStats{
						MuscleGroup: muscleGroup,
					}
				}
				stats := muscleGroupData[muscleGroup]
				stats.TrainingCount++
				// 安全计算总重量：weight * sets * reps
				weight := 0.0
				sets := 0
				reps := 0
				if exercise.Weight != nil {
					weight = *exercise.Weight
				}
				if exercise.Sets != nil {
					sets = *exercise.Sets
				}
				if exercise.Reps != nil {
					reps = *exercise.Reps
				}
				stats.TotalWeight += weight * float64(sets*reps)
				totalCount++
			}
		}
	}

	// Convert map to slice and calculate percentages
	result := []domain.MuscleGroupStats{}
	for _, stats := range muscleGroupData {
		if totalCount > 0 {
			stats.Percentage = (stats.TrainingCount * 100) / totalCount
		}
		result = append(result, *stats)
	}

	// Initialize empty array to avoid null in JSON
	if result == nil {
		result = []domain.MuscleGroupStats{}
	}

	return result, nil
}

func (su *statsUsecase) GetPersonalRecords(c context.Context, userID string) ([]domain.PersonalRecord, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Get all training records for the user
	records, _, err := su.trainingRecordRepository.GetByUserID(ctx, userID, 1, 10000, "", "", 0)
	if err != nil {
		return nil, err
	}

	// Track personal records for each exercise
	prMap := make(map[string]*domain.PersonalRecord)

	for _, record := range records {
		for _, exercise := range record.Exercises {
			if exercise.Name == "" {
				continue
			}

			// Extract date from startTime (YYYY-MM-DD HH:mm:ss -> YYYY-MM-DD)
			recordDate := ""
			if record.StartTime != nil && len(*record.StartTime) >= 10 {
				recordDate = (*record.StartTime)[:10]
			}

			// 获取重量，如果为空则跳过
			weight := 0.0
			if exercise.Weight != nil {
				weight = *exercise.Weight
			}

			// Check if this is a new PR for this exercise
			if pr, exists := prMap[exercise.Name]; exists {
				if weight > pr.MaxWeight {
					pr.MaxWeight = weight
					pr.Date = recordDate
					pr.RecordID = 0 // Would need to store exercise ID to populate this
				}
			} else {
				prMap[exercise.Name] = &domain.PersonalRecord{
					ExerciseName: exercise.Name,
					MaxWeight:    weight,
					Date:         recordDate,
					RecordID:     0,
				}
			}
		}
	}

	// Convert map to slice
	result := []domain.PersonalRecord{}
	for _, pr := range prMap {
		result = append(result, *pr)
	}

	// Initialize empty array to avoid null in JSON
	if result == nil {
		result = []domain.PersonalRecord{}
	}

	return result, nil
}

func (su *statsUsecase) GetCalendar(c context.Context, userID string, year, month int) ([]domain.CalendarDay, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Calculate start and end dates for the month
	startDate := fmt.Sprintf("%d-%02d-01", year, month)

	// Get last day of month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	endDate := lastDay.Format("2006-01-02")

	// Get training records for the month
	records, _, err := su.trainingRecordRepository.GetByUserID(ctx, userID, 1, 1000, startDate, endDate, 0)
	if err != nil {
		return nil, err
	}

	// Create a map to track training by date
	dateMap := make(map[string]*domain.CalendarDay)

	// Initialize all days of the month
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dateMap[dateStr] = &domain.CalendarDay{
			Date:        dateStr,
			HasTraining: false,
		}
	}

	// Fill in training data
	for _, record := range records {
		// Extract date from startTime (YYYY-MM-DD HH:mm:ss -> YYYY-MM-DD)
		recordDate := ""
		if record.StartTime != nil && len(*record.StartTime) >= 10 {
			recordDate = (*record.StartTime)[:10]
		}

		if recordDate == "" {
			continue
		}

		if day, exists := dateMap[recordDate]; exists {
			day.HasTraining = true
			day.TrainingCount++
			if record.Duration != nil {
				day.TotalDuration += *record.Duration
			}
		}
	}

	// Convert map to slice
	result := []domain.CalendarDay{}
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if day, exists := dateMap[dateStr]; exists {
			result = append(result, *day)
		}
	}

	// Initialize empty array to avoid null in JSON
	if result == nil {
		result = []domain.CalendarDay{}
	}

	return result, nil
}

func (su *statsUsecase) GetPlanStats(c context.Context, userID, planID, period string) (domain.PlanStats, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Get the plan
	plan, err := su.fitnessPlanRepository.GetByID(ctx, planID)
	if err != nil {
		return domain.PlanStats{}, err
	}

	// Validate ownership
	if plan.UserID.Hex() != userID {
		return domain.PlanStats{}, fmt.Errorf("unauthorized access to fitness plan")
	}

	// Calculate date range based on period
	now := time.Now()
	var startDate, endDate string

	switch period {
	case "week":
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	case "month":
		startDate = now.AddDate(0, -1, 0).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
	case "whole":
		startDate = plan.StartDate
		endDate = plan.EndDate
	default:
		startDate = now.AddDate(0, 0, -7).Format("2006-01-02")
		endDate = now.Format("2006-01-02")
		period = "week"
	}

	// Get training records for this plan within the period
	// Note: We need to get records related to this plan
	// For now, we use plan stats directly
	totalDays := plan.DurationWeeks * plan.TrainingDaysPerWeek
	completedDays := len(plan.CompletedDays)
	skippedDays := len(plan.SkippedDays)

	// Calculate completion rate
	effectiveTotalDays := totalDays - skippedDays
	completionRate := 0
	if effectiveTotalDays > 0 {
		completionRate = (completedDays * 100) / effectiveTotalDays
	}

	// Build trend data
	trend := []domain.DailyStats{}

	// Parse start date
	planStartDate, err := time.Parse("2006-01-02", plan.StartDate)
	if err == nil {
		// Build daily trend for completed and skipped days
		for i := 0; i < totalDays && i < 30; i++ { // Limit to 30 days
			date := planStartDate.AddDate(0, 0, i).Format("2006-01-02")
			dayNum := i + 1

			// Check if day is completed or skipped
			status := "pending"
			for _, d := range plan.CompletedDays {
				if d == dayNum {
					status = "completed"
					break
				}
			}
			if status == "pending" {
				for _, d := range plan.SkippedDays {
					if d == dayNum {
						status = "skipped"
						break
					}
				}
			}

			trend = append(trend, domain.DailyStats{
				Date:             date,
				CompletionStatus: status,
			})
		}
	}

	stats := domain.PlanStats{
		PlanID:         planID,
		Period:         period,
		StartDate:      startDate,
		EndDate:        endDate,
		CompletionRate: completionRate,
		CompletedDays:  completedDays,
		SkippedDays:    skippedDays,
		TotalDuration:  plan.TotalDuration,
		TotalWeight:    plan.TotalWeight,
		TotalCalories:  plan.TotalCalories,
		Trend:          trend,
	}

	return stats, nil
}

func (su *statsUsecase) GetPlanProgressList(c context.Context, userID, status string, page, pageSize int) ([]domain.PlanProgressSummary, int64, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Get user's plans
	plans, total, err := su.fitnessPlanRepository.GetByUserID(ctx, userID, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Build progress summaries
	result := []domain.PlanProgressSummary{}
	for _, plan := range plans {
		totalDays := plan.DurationWeeks * plan.TrainingDaysPerWeek
		completedDays := len(plan.CompletedDays)
		skippedDays := len(plan.SkippedDays)

		// Calculate completion rate
		effectiveTotalDays := totalDays - skippedDays
		completionRate := 0
		if effectiveTotalDays > 0 {
			completionRate = (completedDays * 100) / effectiveTotalDays
		}

		// Calculate current week and day
		currentWeek := plan.CurrentWeek
		currentDay := plan.CurrentDay
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

		summary := domain.PlanProgressSummary{
			PlanID:         plan.ID.Hex(),
			Name:           plan.Name,
			Status:         plan.Status,
			CompletionRate: completionRate,
			CurrentWeek:    currentWeek,
			CurrentDay:     currentDay,
			EndDate:        plan.EndDate,
		}
		result = append(result, summary)
	}

	// Initialize empty array to avoid null in JSON
	if result == nil {
		result = []domain.PlanProgressSummary{}
	}

	return result, total, nil
}
