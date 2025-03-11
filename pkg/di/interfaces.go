package di

import "linkshorter/internal/user"

type IStatRepository interface {
	AddClick(linkId uint)
}

type IUserReposetory interface {
	CreateUser(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}