package entity

import (
	"html"
	"strings"
	"time"
)

type Note struct {
	ID        uint64    `gorm:"primary_key;auto_increment;" json:"id"`
	Title     string    `gorm:"size:120;not null;" json:"title"`
	Subtitle  string    `gorm:"size:200;" json:"subtitle"`
	Body      string    `gorm:"type:text;" json:"body"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;" json:"updated_at"`
}

type BriefNote struct {
	ID        uint64    `gorm:"primary_key;auto_increment;" json:"id"`
	Title     string    `gorm:"size:120;not null;" json:"title"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;" json:"updated_at"`
}

type Notes []Note

func (note *Note) GetBrief() interface{} {
	return &BriefNote{
		ID:        note.ID,
		Title:     note.Title,
		UpdatedAt: note.UpdatedAt,
	}
}

func (notes Notes) GetBriefs() []interface{} {
	result := make([]interface{}, len(notes))
	for index, note := range notes {
		result[index] = note.GetBrief()
	}
	return result
}

func (note *Note) BeforeSave() {
	note.Title = html.EscapeString(strings.TrimSpace(note.Title))
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
}

func (note *Note) Prepare() {
	note.Title = html.EscapeString(strings.TrimSpace(note.Title))
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
}

func (note *Note) Validate() map[string]string {
	errorMessages := make(map[string]string)

	if note.Title == "" || note.Title == "null" {
		errorMessages["title_required"] = "title required"
	}
	if note.Body == "" || note.Body == "null" {
		errorMessages["body_required"] = "body required"
	}
	return errorMessages
}
