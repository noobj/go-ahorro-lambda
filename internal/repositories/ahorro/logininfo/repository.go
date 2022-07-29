package repository

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	. "github.com/noobj/go-serverless-services/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginInfo struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	User         primitive.ObjectID
	RefreshToken string
	CreatedAt    primitive.Timestamp `bson:"createdAt"`
}

type LoginInfoRepository struct {
	AbstractRepository
}

func New() *LoginInfoRepository {
	abstractRepository := AbstractRepository{
		BaseRepository: BaseRepository{
			Client:     mongodb.GetInstance(),
			Collection: mongodb.GetInstance().Database("ahorro").Collection("loginInfos"),
		},
	}
	repo := LoginInfoRepository{AbstractRepository: abstractRepository}
	repo.IRepository = abstractRepository

	return &repo
}
