package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type planTemplateUsecase struct {
	planTemplateRepository domain.PlanTemplateRepository
	contextTimeout         time.Duration
}

func NewPlanTemplateUsecase(planTemplateRepository domain.PlanTemplateRepository, timeout time.Duration) domain.PlanTemplateUsecase {
	return &planTemplateUsecase{
		planTemplateRepository: planTemplateRepository,
		contextTimeout:         timeout,
	}
}

func (ptu *planTemplateUsecase) GetByID(c context.Context, templateID string) (domain.PlanTemplate, error) {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	return ptu.planTemplateRepository.GetByID(ctx, templateID)
}

func (ptu *planTemplateUsecase) GetList(c context.Context, goal, level, splitType, equipment string, durationWeeksMin, durationWeeksMax, page, pageSize int) ([]domain.PlanTemplate, int64, error) {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	templates, total, err := ptu.planTemplateRepository.GetList(ctx, goal, level, splitType, equipment, durationWeeksMin, durationWeeksMax, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Initialize empty array to avoid null in JSON
	if templates == nil {
		templates = []domain.PlanTemplate{}
	}

	return templates, total, nil
}

func (ptu *planTemplateUsecase) CreateCustom(c context.Context, userID string, request *domain.CreateCustomTemplateRequest) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Initialize arrays to avoid null in JSON
	trainingDays := request.TrainingDays
	if trainingDays == nil {
		trainingDays = []domain.TrainingDay{}
	}

	tags := request.Tags
	if tags == nil {
		tags = []string{}
	}

	template := &domain.PlanTemplate{
		ID:                   primitive.NewObjectID(),
		UserID:               &userObjectID,
		Name:                 request.Name,
		Description:          request.Description,
		Goal:                 request.Goal,
		SplitType:            request.SplitType,
		Level:                request.Level,
		Equipment:            request.Equipment,
		DurationWeeks:        request.DurationWeeks,
		TrainingDaysPerWeek:  request.TrainingDaysPerWeek,
		TrainingDays:         trainingDays,
		Tags:                 tags,
		ImageUrl:             request.ImageUrl,
		RecommendedIntensity: request.RecommendedIntensity,
		Author:               "个人模板",
		IsOfficial:           false,
	}

	err = ptu.planTemplateRepository.Create(ctx, template)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":   template.ID.Hex(),
		"name": template.Name,
	}, nil
}

func (ptu *planTemplateUsecase) CreateOfficial(c context.Context, request *domain.CreateOfficialTemplateRequest) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	// Initialize arrays to avoid null in JSON
	trainingDays := request.TrainingDays
	if trainingDays == nil {
		trainingDays = []domain.TrainingDay{}
	}

	tags := request.Tags
	if tags == nil {
		tags = []string{}
	}

	// 设置默认作者
	author := request.Author
	if author == "" {
		author = "FitEasy官方"
	}

	template := &domain.PlanTemplate{
		ID:                   primitive.NewObjectID(),
		UserID:               nil, // 官方模板没有用户ID
		Name:                 request.Name,
		Description:          request.Description,
		Goal:                 request.Goal,
		SplitType:            request.SplitType,
		Level:                request.Level,
		Equipment:            request.Equipment,
		DurationWeeks:        request.DurationWeeks,
		TrainingDaysPerWeek:  request.TrainingDaysPerWeek,
		TrainingDays:         trainingDays,
		Tags:                 tags,
		ImageUrl:             request.ImageUrl,
		RecommendedIntensity: request.RecommendedIntensity,
		Author:               author,
		IsOfficial:           true, // 标记为官方模板
	}

	err := ptu.planTemplateRepository.Create(ctx, template)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":         template.ID.Hex(),
		"name":       template.Name,
		"isOfficial": true,
	}, nil
}

func (ptu *planTemplateUsecase) Duplicate(c context.Context, userID, templateID string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	// Get original template
	original, err := ptu.planTemplateRepository.GetByID(ctx, templateID)
	if err != nil {
		return nil, errors.New("template not found")
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	// Create a copy of training days
	trainingDays := make([]domain.TrainingDay, len(original.TrainingDays))
	copy(trainingDays, original.TrainingDays)

	// Create a copy of tags
	tags := make([]string, len(original.Tags))
	copy(tags, original.Tags)

	// Create new template as a personal copy
	newTemplate := &domain.PlanTemplate{
		ID:                   primitive.NewObjectID(),
		UserID:               &userObjectID,
		Name:                 "复制自：" + original.Name,
		Description:          original.Description,
		Goal:                 original.Goal,
		SplitType:            original.SplitType,
		Level:                original.Level,
		Equipment:            original.Equipment,
		DurationWeeks:        original.DurationWeeks,
		TrainingDaysPerWeek:  original.TrainingDaysPerWeek,
		TrainingDays:         trainingDays,
		Tags:                 tags,
		ImageUrl:             original.ImageUrl,
		RecommendedIntensity: original.RecommendedIntensity,
		Author:               "个人模板",
		IsOfficial:           false,
	}

	err = ptu.planTemplateRepository.Create(ctx, newTemplate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":   newTemplate.ID.Hex(),
		"name": newTemplate.Name,
	}, nil
}

func (ptu *planTemplateUsecase) Update(c context.Context, userID, templateID string, request *domain.UpdateTemplateRequest) error {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	// Get existing template and validate ownership
	template, err := ptu.planTemplateRepository.GetByID(ctx, templateID)
	if err != nil {
		return errors.New("template not found")
	}

	// Check if user owns this template
	if template.UserID == nil || template.UserID.Hex() != userID {
		return errors.New("unauthorized: you can only update your own templates")
	}

	// Check if it's an official template
	if template.IsOfficial {
		return errors.New("cannot modify official templates")
	}

	// Update fields if provided
	if request.Name != nil {
		template.Name = *request.Name
	}
	if request.Description != nil {
		template.Description = *request.Description
	}
	if request.Goal != nil {
		template.Goal = *request.Goal
	}
	if request.SplitType != nil {
		template.SplitType = *request.SplitType
	}
	if request.Level != nil {
		template.Level = *request.Level
	}
	if request.Equipment != nil {
		template.Equipment = *request.Equipment
	}
	if request.DurationWeeks != nil {
		template.DurationWeeks = *request.DurationWeeks
	}
	if request.TrainingDaysPerWeek != nil {
		template.TrainingDaysPerWeek = *request.TrainingDaysPerWeek
	}
	if request.TrainingDays != nil {
		template.TrainingDays = request.TrainingDays
	}
	if request.Tags != nil {
		template.Tags = request.Tags
	}
	if request.ImageUrl != nil {
		template.ImageUrl = *request.ImageUrl
	}
	if request.RecommendedIntensity != nil {
		template.RecommendedIntensity = *request.RecommendedIntensity
	}

	return ptu.planTemplateRepository.Update(ctx, templateID, &template)
}

func (ptu *planTemplateUsecase) Delete(c context.Context, userID, templateID string) error {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	// Get existing template and validate ownership
	template, err := ptu.planTemplateRepository.GetByID(ctx, templateID)
	if err != nil {
		return errors.New("template not found")
	}

	// Check if user owns this template
	if template.UserID == nil || template.UserID.Hex() != userID {
		return errors.New("unauthorized: you can only delete your own templates")
	}

	// Check if it's an official template
	if template.IsOfficial {
		return errors.New("cannot delete official templates")
	}

	return ptu.planTemplateRepository.Delete(ctx, templateID)
}
