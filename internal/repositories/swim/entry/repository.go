package entry

import (
	"github.com/noobj/swim-crowd-lambda-go/internal/mongodb"
	. "github.com/noobj/swim-crowd-lambda-go/internal/repositories"
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
			Collection: mongodb.GetInstance().Database("ahorro").Collection("entries"),
		},
	}
	repo := EntryRepository{AbstractRepository: abstractRepository}
	repo.IRepository = abstractRepository

	return &repo
}
