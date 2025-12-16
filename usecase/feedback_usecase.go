package usecase

import (
	"context"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
)

type feedbackUsecase struct {
	feedbackRepository domain.FeedbackRepository
	contextTimeout     time.Duration
}

func NewFeedbackUsecase(feedbackRepository domain.FeedbackRepository, timeout time.Duration) domain.FeedbackUsecase {
	return &feedbackUsecase{
		feedbackRepository: feedbackRepository,
		contextTimeout:     timeout,
	}
}

func (fu *feedbackUsecase) Create(c context.Context, feedback *domain.Feedback) error {
	ctx, cancel := context.WithTimeout(c, fu.contextTimeout)
	defer cancel()
	return fu.feedbackRepository.Create(ctx, feedback)
}

