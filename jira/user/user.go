package user

type UserAdapter interface {
	Myself() (User, error)
	User(accountId string) (User, error)
	Users(accountIds []string) ([]User, error)
}

type User struct {
	adapter   UserAdapter
	AccountId string
	Name      string
	Email     string
}

func NewUser(adapter UserAdapter, accountId string, name string, email string) User {
	return User{
		adapter:   adapter,
		AccountId: accountId,
		Name:      name,
		Email:     email,
	}
}
