package user

import "linkshorter/pkg/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Database: database,
	}
}

func (repository *UserRepository) CreateUser(user *User) (*User, error) {
	result := repository.Database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repository *UserRepository) FindByEmail(email string) (*User, error) {
	var user User

	result := repository.Database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}