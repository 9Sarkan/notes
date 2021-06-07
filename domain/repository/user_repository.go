package repository

import "github.com/9sarkan/notes/domain/entity"

type UserRepository interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUser(uint64) (*entity.User, error)
	GetUsers() ([]entity.User, error)
	GetUserByUsernamePassword(*entity.User) (*entity.User, map[string]string)
}
