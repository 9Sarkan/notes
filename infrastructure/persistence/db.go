package persistence

import (
	"fmt"

	"github.com/9sarkan/notes/domain/entity"
	"github.com/9sarkan/notes/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repositories struct {
	Note repository.NoteRepository
	User repository.UserRepository
	db   *gorm.DB
}

func NewRepositories(DbDriver, DbUser, DbPassword, DbHost, DbPort, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=false password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(DbDriver, DBURL)
	if err != nil {
		return nil, err
	}
	return &Repositories{
		User: NewUserRepository(db),
		Note: NewNoteRepository(db),
		db:   db,
	}, nil
}

func (r *Repositories) Close() error {
	return r.db.Close()
}
func (r *Repositories) Automigrate() error {
	return r.db.AutoMigrate(&entity.Note{}, &entity.User{}).Error
}
