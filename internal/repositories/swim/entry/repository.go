package entry

import (
	"github.com/noobj/go-serverless-services/internal/mongodb"
	. "github.com/noobj/go-serverless-services/internal/repositories"
)

type Entry struct {
	Amount int    `json:"amount"`
	Time   string `json:"time"`
}

type EntryRepository struct {
	AbstractRepository
}

func New() *EntryRepository {
	abstractRepository := AbstractRepository{
		BaseRepository: BaseRepository{
			Client:     mongodb.GetInstance(),
			Collection: mongodb.GetInstance().Database("swimCrowdDB").Collection("entries"),
		},
	}
	repo := EntryRepository{AbstractRepository: abstractRepository}
	repo.IRepository = abstractRepository

	return &repo
}
