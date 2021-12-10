package models

import (
	"gorm.io/gorm"
)

// Template tempalate
type Template struct {
	ID        string
	Name      string
	Title     string
	Content   string
	CreateBy  string
	Status    int
	CreatedAt int64
	UpdatedAt int64
}

// TemplateRepo TemplateRepo
type TemplateRepo interface {
	Create(*gorm.DB, *Template) error
	UpdateTemplate(*gorm.DB, *Template) error
	Get(*gorm.DB, string) (*Template, error)
	Delete(*gorm.DB, string) error
	QueryTemplate(*gorm.DB, string, int, int) ([]*Template, int64, error)
}
