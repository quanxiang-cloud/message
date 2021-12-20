package mysql

import (
	"git.internal.yunify.com/qxp/misc/time2"
	"github.com/quanxiang-cloud/message/internal/models"
	"gorm.io/gorm"
)

type tempalateRepo struct {
}

func (t *tempalateRepo) TableName() string {
	return "template"

}

// NewTemplateRepo create
func NewTemplateRepo() models.TemplateRepo {
	return &tempalateRepo{}
}

// Create Create
func (t *tempalateRepo) Create(db *gorm.DB, template *models.Template) error {
	return db.Table(t.TableName()).Create(template).Error
}

// UpdateTemplate update
func (t *tempalateRepo) UpdateTemplate(db *gorm.DB, template *models.Template) error {
	return db.Table(t.TableName()).
		Where("id = ?", template.ID).Updates(map[string]interface{}{
		"content":    template.Content,
		"title":      template.Title,
		"name":       template.Name,
		"updated_at": time2.NowUnix(),
	}).Error
}

func (t *tempalateRepo) Get(db *gorm.DB, id string) (*models.Template, error) {
	template := new(models.Template)
	err := db.Table(t.TableName()).
		Where("id = ?", id).
		Find(template).Error
	if err != nil {
		return nil, err
	}
	if template.ID == "" {
		return nil, nil
	}
	return template, nil
}

func (t *tempalateRepo) Delete(db *gorm.DB, templateID string) error {
	return db.Table(t.TableName()).Where("id = ?", templateID).Delete(&models.Template{}).Error
}

func (t *tempalateRepo) QueryTemplate(db *gorm.DB, title string, page, limit int) ([]*models.Template, int64, error) {

	ql := db.Table(t.TableName())

	if title != "" {
		ql = ql.Where("title like ", title)
	}
	var total int64
	ql.Count(&total)
	ql = ql.Limit(limit).Offset((page - 1) * limit)
	ql = ql.Order("created_at desc")
	templates := make([]*models.Template, 0)
	err := ql.Find(&templates).Error
	return templates, total, err
}
