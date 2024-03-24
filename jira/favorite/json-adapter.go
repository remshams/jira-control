package favorite

import (
	"encoding/json"

	"github.com/charmbracelet/log"
	file_store "github.com/remshams/common/utils/file-store"
)

type FavoriteDto = Favorite
type FavoriteDtoList = []FavoriteDto

func replaceOrAddFavoriteInList(favorites []Favorite, newFavorite Favorite) FavoriteDtoList {
	for i, favoriteDto := range favorites {
		if favoriteDto.Id == newFavorite.Id {
			favorites[i] = newFavorite
			return favorites
		}
	}
	return append(favorites, newFavorite)
}

func createFavoriteJson(favorites []FavoriteDto) ([]byte, error) {
	favoritesJson, err := json.Marshal(favorites)
	if err != nil {
		log.Errorf("FavoriteDto: Could not marshal favorites: %v", err)
		return nil, err
	}
	return favoritesJson, nil
}

func favoritesFromJson(data []byte) ([]Favorite, error) {
	var favorites []Favorite
	err := json.Unmarshal(data, &favorites)
	if err != nil {
		log.Errorf("FavoriteDto: Could not unmarshal favorites: %v", err)
		return nil, err
	}
	return favorites, nil
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
	favoriteJson, err := createFavoriteJson(favorites)
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
	favorites, err := favoritesFromJson(data)
	if err != nil {
		return nil, err
	}
	return favorites, nil
}
