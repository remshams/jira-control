package user

type UserAdapter interface {
	Myself() (User, error)
	User(accountId string) (User, error)
}

type User struct {
	AccountId string
	Name      string
	Email     string
}

func NewUser(accountId string, name string, email string) User {
	return User{
		AccountId: accountId,
		Name:      name,
		Email:     email,
	}
}
