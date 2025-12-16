package repository

import (
	"context"
	"time"

	"github.com/zhengshui/flow-link-server/domain"
	"github.com/zhengshui/flow-link-server/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type feedbackRepository struct {
	database   mongo.Database
	collection string
}

func NewFeedbackRepository(db mongo.Database, collection string) domain.FeedbackRepository {
	return &feedbackRepository{
		database:   db,
		collection: collection,
	}
}

func (fr *feedbackRepository) Create(c context.Context, feedback *domain.Feedback) error {
	collection := fr.database.Collection(fr.collection)
	feedback.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := collection.InsertOne(c, feedback)
	return err
}

