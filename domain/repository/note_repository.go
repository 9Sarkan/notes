package repository

import "github.com/9sarkan/notes/domain/entity"

type NoteRepository interface {
	GetNote(uint64) (*entity.Note, error)
	GetNotes() ([]entity.Note, error)
	SaveNote(*entity.Note) (*entity.Note, map[string]string)
	UpdateNote(*entity.Note) (*entity.Note, map[string]string)
	DeleteNote(uint64) error
}
