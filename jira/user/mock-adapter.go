package user

type MockUserAdapter struct {
}

func NewMockUserAdapter() MockUserAdapter {
	return MockUserAdapter{}
}

func (_ MockUserAdapter) Myself() (User, error) {
	return NewUser("1", "mock user", "mock.user@mock.com"), nil
}
