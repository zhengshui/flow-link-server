package usecase

import (
	"context"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
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

func (ptu *planTemplateUsecase) GetList(c context.Context, goal, level string, page, pageSize int) ([]domain.PlanTemplate, int64, error) {
	ctx, cancel := context.WithTimeout(c, ptu.contextTimeout)
	defer cancel()

	templates, total, err := ptu.planTemplateRepository.GetList(ctx, goal, level, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Initialize empty array to avoid null in JSON
	if templates == nil {
		templates = []domain.PlanTemplate{}
	}

	return templates, total, nil
}
