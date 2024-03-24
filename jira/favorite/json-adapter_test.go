package favorite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var favorite = NewFavorite(nil, "NX-Testing", 9)

func createTestFavorite() Favorite {
	return NewFavorite(nil, "NX-Testing", 9)
}

func TestReplaceOrAddFavoriteInList_Existing(t *testing.T) {
	existingFavorite := createTestFavorite()
	newFavorite := createTestFavorite()
	newFavorite.Id = existingFavorite.Id
	newHoursSpent := 10.5
	newFavorite.HoursSpent = newHoursSpent

	favorites := []Favorite{existingFavorite}
	updatedFavorites := replaceOrAddFavoriteInList(favorites, newFavorite)

	assert.Equal(t, 1, len(updatedFavorites))
	assert.Equal(t, newHoursSpent, updatedFavorites[0].HoursSpent)
}

func TestReplaceOrAddFavoriteInList_NotExisting(t *testing.T) {
	existingFavorite := createTestFavorite()
	newFavorite := createTestFavorite()

	favorites := []Favorite{existingFavorite}
	updatedFavorites := replaceOrAddFavoriteInList(favorites, newFavorite)

	assert.Equal(t, 2, len(updatedFavorites))
}

func TestReplaceOrAddFavoriteInList_Empty(t *testing.T) {
	newFavorite := createTestFavorite()

	updatedFavorites := replaceOrAddFavoriteInList([]Favorite{}, newFavorite)

	assert.Equal(t, 1, len(updatedFavorites))
}

func TestSortByLastUsed_Ordered(t *testing.T) {
	favoriteFirst := createTestFavorite()
	favoriteSecond := createTestFavorite()
	favoriteFirst.LastUsedAt = favoriteSecond.LastUsedAt.AddDate(0, 0, 1)
	favoriteList := []Favorite{favoriteSecond, favoriteFirst}
	sortedFavorites := sortByLastUsedAt(favoriteList)

	assert.Equal(t, favoriteFirst, sortedFavorites[0])
	assert.Equal(t, favoriteSecond, sortedFavorites[1])
}

func TestSortByLastUsed_Unordered(t *testing.T) {
	favoriteFirst := createTestFavorite()
	favoriteSecond := createTestFavorite()
	favoriteFirst.LastUsedAt = favoriteSecond.LastUsedAt.AddDate(0, 0, 1)
	favoriteList := []Favorite{favoriteSecond, favoriteFirst}
	sortedFavorites := sortByLastUsedAt(favoriteList)

	assert.Equal(t, favoriteFirst, sortedFavorites[0])
	assert.Equal(t, favoriteSecond, sortedFavorites[1])
}

func TestSortByLastUsed_Empty(t *testing.T) {
	sortedFavorites := sortByLastUsedAt([]Favorite{})

	assert.Equal(t, 0, len(sortedFavorites))
}
