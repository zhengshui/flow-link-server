package domain

import "context"

// DailyStats 每日统计
type DailyStats struct {
	Date             string  `json:"date"`                       // 日期 YYYY-MM-DD
	TrainingCount    int     `json:"trainingCount"`              // 当日训练次数
	Duration         int     `json:"duration"`                   // 当日训练时长
	Weight           float64 `json:"weight"`                     // 当日总重量
	Sets             int     `json:"sets"`                       // 当日总组数
	Calories         int     `json:"calories"`                   // 当日消耗卡路里
	CompletionStatus string  `json:"completionStatus,omitempty"` // 完成状态(用于计划统计)
}

// TrainingStats 训练统计
type TrainingStats struct {
	UserID              string       `json:"userId"`
	Period              string       `json:"period"`              // 统计周期(week/month/year)
	StartDate           string       `json:"startDate"`           // 统计开始日期
	EndDate             string       `json:"endDate"`             // 统计结束日期
	TotalTrainingCount  int          `json:"totalTrainingCount"`  // 总训练次数
	TotalDuration       int          `json:"totalDuration"`       // 总训练时长(分钟)
	TotalWeight         float64      `json:"totalWeight"`         // 总重量(kg)
	TotalSets           int          `json:"totalSets"`           // 总组数
	TotalCalories       int          `json:"totalCalories"`       // 总消耗卡路里
	AvgDuration         int          `json:"avgDuration"`         // 平均训练时长
	AvgWeight           float64      `json:"avgWeight"`           // 平均单次重量
	MostTrainedMuscle   string       `json:"mostTrainedMuscle"`   // 训练最多的肌群
	FavoriteExercise    string       `json:"favoriteExercise"`    // 最常做的训练项目
	DailyStats          []DailyStats `json:"dailyStats"`          // 每日统计数据
}

// MuscleGroupStats 肌群训练统计
type MuscleGroupStats struct {
	MuscleGroup   string  `json:"muscleGroup"`   // 肌群名称
	TrainingCount int     `json:"trainingCount"` // 训练次数
	TotalWeight   float64 `json:"totalWeight"`   // 总重量
	Percentage    int     `json:"percentage"`    // 占比百分比
}

// PersonalRecord 个人记录
type PersonalRecord struct {
	ExerciseName string  `json:"exerciseName"` // 训练项目名称
	MaxWeight    float64 `json:"maxWeight"`    // 最大重量
	Date         string  `json:"date"`         // 创建日期
	RecordID     int     `json:"recordId"`     // 记录ID
}

// CalendarDay 日历天数据
type CalendarDay struct {
	Date          string `json:"date"`          // 日期 YYYY-MM-DD
	HasTraining   bool   `json:"hasTraining"`   // 是否有训练
	TrainingCount int    `json:"trainingCount"` // 训练次数
	TotalDuration int    `json:"totalDuration"` // 总时长
}

// PlanStats 计划维度统计
type PlanStats struct {
	PlanID         string       `json:"planId"`
	Period         string       `json:"period"`         // 统计周期(week/month/whole)
	StartDate      string       `json:"startDate"`      // 统计开始日期
	EndDate        string       `json:"endDate"`        // 统计结束日期
	CompletionRate int          `json:"completionRate"` // 完成率
	CompletedDays  int          `json:"completedDays"`  // 完成天数
	SkippedDays    int          `json:"skippedDays"`    // 跳过天数
	TotalDuration  int          `json:"totalDuration"`  // 累计训练时长
	TotalWeight    float64      `json:"totalWeight"`    // 累计重量
	TotalCalories  int          `json:"totalCalories"`  // 累计消耗卡路里
	Trend          []DailyStats `json:"trend"`          // 趋势数据
}

// PlanProgressSummary 计划进度概览
type PlanProgressSummary struct {
	PlanID         string `json:"planId"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	CompletionRate int    `json:"completionRate"`
	CurrentWeek    int    `json:"currentWeek"`
	CurrentDay     int    `json:"currentDay"`
	EndDate        string `json:"endDate"`
}

// StatsUsecase 统计用例接口
type StatsUsecase interface {
	GetTrainingStats(c context.Context, userID string, period, startDate, endDate string) (TrainingStats, error)
	GetMuscleGroupStats(c context.Context, userID string, period string) ([]MuscleGroupStats, error)
	GetPersonalRecords(c context.Context, userID string) ([]PersonalRecord, error)
	GetCalendar(c context.Context, userID string, year, month int) ([]CalendarDay, error)
	GetPlanStats(c context.Context, userID, planID, period string) (PlanStats, error)
	GetPlanProgressList(c context.Context, userID, status string, page, pageSize int) ([]PlanProgressSummary, int64, error)
}
