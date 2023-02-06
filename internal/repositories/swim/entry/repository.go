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
	IRepository
}

func New() *EntryRepository {
	baseRepository := BaseRepository{
		Collection: mongodb.GetInstance().Database("swimCrowdDB").Collection("entries"),
	}
	repo := EntryRepository{IRepository: baseRepository}

	return &repo
}
