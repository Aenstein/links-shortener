package auth_test

import (
	"linkshorter/internal/auth"
	"linkshorter/internal/user"
	"testing"
)

type MockUserRepository struct {}

func (r *MockUserRepository) CreateUser(user *user.User) (*user.User, error) {
	return user, nil
}

func (r *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "a@a.ru"
	authService := auth.NewAuthService(&MockUserRepository{})

	email, err := authService.Register(initialEmail, "1", "Vanya")
	if err != nil {
		t.Fatal(err)
	}

	if email != initialEmail {
		t.Fatalf("Email %s do not match %s", email, initialEmail)
	}
}