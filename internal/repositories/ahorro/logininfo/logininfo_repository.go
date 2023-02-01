package repository

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	. "github.com/noobj/go-serverless-services/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginInfo struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	User         primitive.ObjectID
	RefreshToken string             `bson:"refreshToken"`
	CreatedAt    primitive.DateTime `bson:"createdAt"`
}

type LoginInfoRepository struct {
	IRepository
}

func New() *LoginInfoRepository {
	baseRepository := BaseRepository{
		Client:     mongodb.GetInstance(),
		Collection: mongodb.GetInstance().Database("ahorro").Collection("loginInfos"),
	}
	repo := LoginInfoRepository{IRepository: baseRepository}

	return &repo
}
