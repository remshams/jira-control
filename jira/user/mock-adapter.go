package user

import (
	"fmt"

	"github.com/charmbracelet/log"
)

type MockUserAdapter struct {
}

func NewMockUserAdapter() MockUserAdapter {
	return MockUserAdapter{}
}

func (_ MockUserAdapter) Myself() (User, error) {
	log.Debug("MockUserAdapter: Request myself")
	return NewUser("1", "mock user", "mock.user@mock.com"), nil
}

func (_ MockUserAdapter) User(accountId string) (User, error) {
	log.Debugf("MockUserAdapter: Requesting user with accountId: %s", accountId)
	return NewUser("2", fmt.Sprintf("User with accountId: %s", accountId), fmt.Sprintf("mock.%s@mock.com", accountId)), nil
}

func (_ MockUserAdapter) Users(accountIds []string) ([]User, error) {
	log.Debugf("MockUserAdapter: Requesting user with #accountIds: %d", len(accountIds))
	users := []User{}
	for _, accountId := range accountIds {
		users = append(
			users,
			NewUser("2", fmt.Sprintf("User with accountId: %s", accountId), fmt.Sprintf("mock.%s@mock.com", accountId)),
		)
	}
	return users, nil
}
