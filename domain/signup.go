package domain

import (
	"context"
)

type SignupRequest struct {
	Username     string  `json:"username" binding:"required,min=4,max=20"`
	Password     string  `json:"password" binding:"required,min=6,max=20"`
	Nickname     string  `json:"nickname"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	Gender       string  `json:"gender"`        // 男/女
	Age          int     `json:"age"`
	Height       float64 `json:"height"`        // 身高cm
	Weight       float64 `json:"weight"`        // 体重kg
	TargetWeight float64 `json:"targetWeight"`  // 目标体重kg
	FitnessGoal  string  `json:"fitnessGoal"`   // 增肌/减脂/力量提升/耐力提升/综合健身
}

type SignupResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type SignupUsecase interface {
	Create(c context.Context, user *User) error
	GetUserByUsername(c context.Context, username string) (User, error)
	CreateAccessToken(user *User, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user *User, secret string, expiry int) (refreshToken string, err error)
}
