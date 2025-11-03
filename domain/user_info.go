package domain

import "context"

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	Nickname     string  `json:"nickname,omitempty"`
	AvatarUrl    string  `json:"avatarUrl,omitempty"`
	Email        string  `json:"email,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Gender       string  `json:"gender,omitempty"`
	Age          int     `json:"age,omitempty"`
	Height       float64 `json:"height,omitempty"`
	Weight       float64 `json:"weight,omitempty"`
	TargetWeight float64 `json:"targetWeight,omitempty"`
	FitnessGoal  string  `json:"fitnessGoal,omitempty"`
	JoinDate     string  `json:"joinDate"`
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	Nickname     string  `json:"nickname"`
	AvatarUrl    string  `json:"avatarUrl"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	Gender       string  `json:"gender"`
	Age          int     `json:"age"`
	Height       float64 `json:"height"`
	Weight       float64 `json:"weight"`
	TargetWeight float64 `json:"targetWeight"`
	FitnessGoal  string  `json:"fitnessGoal"`
}

// UserInfoUsecase 用户信息用例接口
type UserInfoUsecase interface {
	GetUserInfo(c context.Context, userID string) (User, error)
	UpdateUserInfo(c context.Context, userID string, request *UpdateUserInfoRequest) error
}
