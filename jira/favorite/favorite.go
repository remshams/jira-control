package favorite

import (
	"time"

	"github.com/google/uuid"
)

type FavoriteAdapter interface {
	Store(favorite Favorite) error
	Load() ([]Favorite, error)
}

type Favorite struct {
	adapter    FavoriteAdapter
	Id         uuid.UUID
	IssueKey   string
	HoursSpent float64
	Comment    string
	LastUsedAt time.Time
	CreatedAt  time.Time
}

func NewFavorite(adapter FavoriteAdapter, issueKey string, hoursSpent float64) Favorite {
	return Favorite{
		adapter:    adapter,
		Id:         uuid.New(),
		IssueKey:   issueKey,
		HoursSpent: hoursSpent,
		Comment:    "",
		LastUsedAt: time.Now(),
		CreatedAt:  time.Now(),
	}
}

func (f Favorite) Store() error {
	return f.adapter.Store(f)
}

func (f Favorite) IsEqual(compare Favorite) bool {
	return f.IssueKey == compare.IssueKey && f.HoursSpent == compare.HoursSpent && f.Comment == compare.Comment
}

func (f Favorite) DoesFavoriteAlreadyExist(favorites []Favorite) bool {
	for _, favorite := range favorites {
		if f.IsEqual(favorite) {
			return true
		}
	}
	return false
}
