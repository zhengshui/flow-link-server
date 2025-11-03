package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Password     string             `bson:"password" json:"-"` // 不在JSON中返回密码
	Nickname     string             `bson:"nickname" json:"nickname,omitempty"`
	AvatarUrl    string             `bson:"avatarUrl" json:"avatarUrl,omitempty"`
	Email        string             `bson:"email" json:"email,omitempty"`
	Phone        string             `bson:"phone" json:"phone,omitempty"`
	Gender       string             `bson:"gender" json:"gender,omitempty"`       // 男/女
	Age          int                `bson:"age" json:"age,omitempty"`
	Height       float64            `bson:"height" json:"height,omitempty"`       // 身高(cm)
	Weight       float64            `bson:"weight" json:"weight,omitempty"`       // 体重(kg)
	TargetWeight float64            `bson:"targetWeight" json:"targetWeight,omitempty"` // 目标体重(kg)
	FitnessGoal  string             `bson:"fitnessGoal" json:"fitnessGoal,omitempty"`   // 健身目标
	Role         string             `bson:"role" json:"role"`                     // user/admin
	JoinDate     string             `bson:"joinDate" json:"joinDate"`             // 加入日期 YYYY-MM-DD
	CreatedAt    primitive.DateTime `bson:"createdAt" json:"-"`
	UpdatedAt    primitive.DateTime `bson:"updatedAt" json:"-"`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	Fetch(c context.Context) ([]User, error)
	GetByEmail(c context.Context, email string) (User, error)
	GetByUsername(c context.Context, username string) (User, error)
	GetByID(c context.Context, id string) (User, error)
	Update(c context.Context, id string, user *User) error
}
