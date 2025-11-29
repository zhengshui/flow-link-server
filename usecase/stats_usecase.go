package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
)

type statsUsecase struct {
	trainingRecordRepository domain.TrainingRecordRepository
	contextTimeout           time.Duration
}

func NewStatsUsecase(trainingRecordRepository domain.TrainingRecordRepository, timeout time.Duration) domain.StatsUsecase {
	return &statsUsecase{
		trainingRecordRepository: trainingRecordRepository,
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
		totalDuration += record.Duration
		totalWeight += record.TotalWeight
		totalSets += record.TotalSets
		totalCalories += record.CaloriesBurned

		// Extract date from startTime (YYYY-MM-DD HH:mm:ss -> YYYY-MM-DD)
		recordDate := record.StartTime
		if len(record.StartTime) >= 10 {
			recordDate = record.StartTime[:10]
		}

		// Track daily stats
		if _, exists := dailyStatsMap[recordDate]; !exists {
			dailyStatsMap[recordDate] = &domain.DailyStats{
				Date: recordDate,
			}
		}
		dailyStats := dailyStatsMap[recordDate]
		dailyStats.TrainingCount++
		dailyStats.Duration += record.Duration
		dailyStats.Weight += record.TotalWeight
		dailyStats.Sets += record.TotalSets
		dailyStats.Calories += record.CaloriesBurned

		// Track muscle groups and exercises
		for _, exercise := range record.Exercises {
			if exercise.MuscleGroup != "" {
				muscleGroupCount[exercise.MuscleGroup]++
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
			if exercise.MuscleGroup != "" {
				if _, exists := muscleGroupData[exercise.MuscleGroup]; !exists {
					muscleGroupData[exercise.MuscleGroup] = &domain.MuscleGroupStats{
						MuscleGroup: exercise.MuscleGroup,
					}
				}
				stats := muscleGroupData[exercise.MuscleGroup]
				stats.TrainingCount++
				stats.TotalWeight += exercise.Weight * float64(exercise.Sets*exercise.Reps)
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
			recordDate := record.StartTime
			if len(record.StartTime) >= 10 {
				recordDate = record.StartTime[:10]
			}

			// Check if this is a new PR for this exercise
			if pr, exists := prMap[exercise.Name]; exists {
				if exercise.Weight > pr.MaxWeight {
					pr.MaxWeight = exercise.Weight
					pr.Date = recordDate
					pr.RecordID = 0 // Would need to store exercise ID to populate this
				}
			} else {
				prMap[exercise.Name] = &domain.PersonalRecord{
					ExerciseName: exercise.Name,
					MaxWeight:    exercise.Weight,
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
		recordDate := record.StartTime
		if len(record.StartTime) >= 10 {
			recordDate = record.StartTime[:10]
		}

		if day, exists := dateMap[recordDate]; exists {
			day.HasTraining = true
			day.TrainingCount++
			day.TotalDuration += record.Duration
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
