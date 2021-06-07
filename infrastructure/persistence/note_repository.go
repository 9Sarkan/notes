package persistence

import (
	"errors"

	"github.com/9sarkan/notes/domain/entity"
	"github.com/9sarkan/notes/domain/repository"
	"github.com/jinzhu/gorm"
)

type NoteRepo struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) *NoteRepo {
	return &NoteRepo{
		db: db,
	}
}

var _ repository.NoteRepository = &NoteRepo{}

func (r *NoteRepo) GetNote(id uint64) (*entity.Note, error) {
	var note entity.Note
	err := r.db.Debug().Where("id = ?", id).Find(&note).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("note not found")
	} else if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepo) GetNotes() ([]entity.Note, map[string]string) {
	var notes []entity.Note
	errDB := make(map[string]string)
	err := r.db.Debug().Find(&notes).Error
	if gorm.IsRecordNotFoundError(err) {
		errDB["no_note"] = "notes not found"
		return nil, errDB
	} else if err != nil {
		errDB["db_error"] = "some db error happend"
		return nil, errDB
	}
	return notes, nil
}

func (r *NoteRepo) SaveNote(n *entity.Note) (*entity.Note, map[string]string) {
	errDB := make(map[string]string)
	err := r.db.Debug().Create(&n).Error
	if err != nil {
		errDB["db_error"] = "some error on database happend"
		return nil, errDB
	}
	return n, nil
}
func (r *NoteRepo) UpdateNote(n *entity.Note) (*entity.Note, map[string]string) {
	errDB := make(map[string]string)
	err := r.db.Debug().Save(&n).Error
	if err != nil {
		errDB["db_error"] = "database error happend"
		return nil, errDB
	}
	return n, nil
}
func (r *NoteRepo) DeleteNote(id uint64) error {
	var note entity.Note
	err := r.db.Debug().Where("id = ?", id).Delete(&note).Error
	if err != nil {
		return err
	}
	return nil
}
