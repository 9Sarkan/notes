package persistence

import (
	"errors"
	"strings"

	"github.com/9sarkan/notes/domain/entity"
	"github.com/9sarkan/notes/domain/repository"
	"github.com/9sarkan/notes/infrastructure/security"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

var _ repository.UserRepository = &UserRepo{}

func (r *UserRepo) GetUser(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.Debug().Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepo) GetUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Debug().Find(&users).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (r *UserRepo) GetUserByUsernamePassword(u *entity.User) (*entity.User, map[string]string) {
	var user entity.User
	dbErr := make(map[string]string)
	err := r.db.Debug().Where("username = ?", u.Username).Find(&user).Error
	if gorm.IsRecordNotFoundError(err) {
		dbErr["no_user"] = "user not found"
		return nil, dbErr
	}
	if err != nil {
		dbErr["db_error"] = "database error happend"
		return nil, dbErr
	}
	err = security.CompareHashAndPassword(user.Password, u.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		dbErr["invalid_password"] = "password is invalid"
		return nil, dbErr
	}
	return &user, nil
}

func (r *UserRepo) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	errDB := make(map[string]string)
	err := r.db.Debug().Create(&user).Error
	if err != nil {
		// if username taken
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			errDB["duplicate_user"] = "username has been taken"
			return nil, errDB
		}
		// any other error
		errDB["db_error"] = "some db error happend"
		return nil, errDB
	}
	return user, nil
}
