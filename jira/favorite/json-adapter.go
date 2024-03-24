package favorite

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	file_store "github.com/remshams/common/utils/file-store"
)

type FavoriteDto struct {
	Id         string    `json:"id"`
	IssueKey   string    `json:"issueKey"`
	HoursSpent float64   `json:"hoursSpent"`
	Comment    string    `json:"comment"`
	LastUsedAt time.Time `json:"lastUsedAt"`
	CreatedAt  time.Time `json:"createdAt"`
}

func fromFavorite(favorite Favorite) FavoriteDto {
	return FavoriteDto{
		Id:         favorite.Id.String(),
		IssueKey:   favorite.IssueKey,
		HoursSpent: favorite.HoursSpent,
		Comment:    favorite.Comment,
		LastUsedAt: favorite.LastUsedAt,
		CreatedAt:  favorite.CreatedAt,
	}
}

func (f FavoriteDto) toFavorite(adapter FavoriteAdapter) (Favorite, error) {
	id, err := uuid.Parse(f.Id)
	if err != nil {
		log.Errorf("FavoriteDto: Could not parse id: %v", err)
		return Favorite{}, err
	}
	return Favorite{
		adapter:    adapter,
		Id:         id,
		IssueKey:   f.IssueKey,
		HoursSpent: f.HoursSpent,
		Comment:    f.Comment,
		LastUsedAt: f.LastUsedAt,
		CreatedAt:  f.CreatedAt,
	}, nil
}

func favoritesToJson(favorites []Favorite) ([]byte, error) {
	var favoriteDtos []FavoriteDto
	for _, favorite := range favorites {
		favoriteDtos = append(favoriteDtos, fromFavorite(favorite))
	}
	favoritesJson, err := json.Marshal(favoriteDtos)
	if err != nil {
		log.Errorf("FavoriteDto: Could not marshal favorites: %v", err)
		return nil, err
	}
	return favoritesJson, nil
}

func favoritesFromJson(data []byte, adapter FavoriteAdapter) ([]Favorite, error) {
	var favoriteDtos []FavoriteDto
	err := json.Unmarshal(data, &favoriteDtos)
	if err != nil {
		log.Errorf("FavoriteDto: Could not unmarshal favorites: %v", err)
		return nil, err
	}
	var favorites []Favorite
	for _, favoriteDto := range favoriteDtos {
		favorite, err := favoriteDto.toFavorite(adapter)
		if err != nil {
			continue
		}
		favorites = append(favorites, favorite)
	}
	return favorites, nil
}

func replaceOrAddFavoriteInList(favorites []Favorite, newFavorite Favorite) []Favorite {
	for i, favoriteDto := range favorites {
		if favoriteDto.Id == newFavorite.Id {
			favorites[i] = newFavorite
			return favorites
		}
	}
	return append(favorites, newFavorite)
}

func sortByLastUsedAt(favorites []Favorite) []Favorite {
	sort.Slice(favorites, func(i, j int) bool {
		return favorites[i].LastUsedAt.After(favorites[j].LastUsedAt)
	})
	return favorites
}

func createFavoriteJson(favorites []FavoriteDto) ([]byte, error) {
	favoritesJson, err := json.Marshal(favorites)
	if err != nil {
		log.Errorf("FavoriteDto: Could not marshal favorites: %v", err)
		return nil, err
	}
	return favoritesJson, nil
}

type FavoriteJsonAdapter struct {
	path string
}

func NewFavoriteJsonAdapter(path string) FavoriteJsonAdapter {
	return FavoriteJsonAdapter{
		path: path,
	}
}

func (f FavoriteJsonAdapter) Store(favorite Favorite) error {
	log.Debugf("FavoriteJsonAdapter: Storing favorite %v", favorite)
	favorites, err := f.Load()
	log.Debugf("FavoriteJsonAdapter: favorites %v", favorites)
	if err != nil {
		return err
	}
	favorites = replaceOrAddFavoriteInList(favorites, favorite)
	favorites = sortByLastUsedAt(favorites)
	favoriteJson, err := favoritesToJson(favorites)
	if err != nil {
		return err
	}
	err = file_store.CreateOrUpdateFile(f.path, favoriteJson)
	if err != nil {
		log.Errorf("FavoriteJsonAdapter: Could not store favorite: %v", err)
	}
	return err
}

func (f FavoriteJsonAdapter) Load() ([]Favorite, error) {
	log.Debugf("FavoriteJsonAdapter: Loading favorites from %v", f.path)
	data := file_store.LoadFile(f.path)
	if len(data) == 0 {
		return []Favorite{}, nil
	}
	favorites, err := favoritesFromJson(data, f)
	if err != nil {
		return nil, err
	}
	return favorites, nil
}
