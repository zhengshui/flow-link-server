package usecase

import (
	"context"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userInfoUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewUserInfoUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserInfoUsecase {
	return &userInfoUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *userInfoUsecase) GetUserInfo(c context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()
	return uu.userRepository.GetByID(ctx, userID)
}

func (uu *userInfoUsecase) UpdateUserInfo(c context.Context, userID string, request *domain.UpdateUserInfoRequest) error {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	// 获取当前用户信息
	user, err := uu.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 更新字段
	if request.Nickname != "" {
		user.Nickname = request.Nickname
	}
	if request.AvatarUrl != "" {
		user.AvatarUrl = request.AvatarUrl
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	if request.Gender != "" {
		user.Gender = request.Gender
	}
	if request.Age > 0 {
		user.Age = request.Age
	}
	if request.Height > 0 {
		user.Height = request.Height
	}
	if request.Weight > 0 {
		user.Weight = request.Weight
	}
	if request.TargetWeight > 0 {
		user.TargetWeight = request.TargetWeight
	}
	if request.FitnessGoal != "" {
		user.FitnessGoal = request.FitnessGoal
	}

	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	return uu.userRepository.Update(ctx, userID, &user)
}
