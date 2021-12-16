package mysql

import (
	"github.com/quanxiang-cloud/message/internal/models"
	"gorm.io/gorm"
)

type messageRepo struct {
}

// TableName TableName
func (m *messageRepo) TableName() string {
	return "message_list"
}

// NewMessageRepo createMessageRepo
func NewMessageRepo() models.MessageRepo {
	return &messageRepo{}
}

// Create create
func (m *messageRepo) Create(db *gorm.DB, message *models.MessageList) error {
	return db.Table(m.TableName()).Create(message).Error
}

// Get Get
func (m *messageRepo) Get(db *gorm.DB, id string) (*models.MessageList, error) {
	message := new(models.MessageList)
	err := db.Table(m.TableName()).Where("id = ?", id).Find(message).Error
	if err != nil {
		return nil, err
	}
	return message, nil

}

// Delete Delete
func (m *messageRepo) Delete(db *gorm.DB, id string) error {
	// 查询id
	return db.Table(m.TableName()).Where("id = ?", id).
		Delete(&models.MessageList{}).
		Error
}

// UpdateCount UpdateCount
func (m *messageRepo) UpdateCount(db *gorm.DB, message *models.MessageList) error {
	return db.Table(m.TableName()).Where("id = ?", message.ID).Updates(
		map[string]interface{}{
			"status":   message.Status,
			"send_num": message.SendNum,
			"success":  message.Success,
			"fail":     message.Fail,
		}).Error
}

// List List
func (m *messageRepo) List(db *gorm.DB, status, types int8, keyword string, page, limit int) ([]*models.MessageList, int64, error) {
	ql := db.Table(m.TableName())
	if status != 0 {
		ql = ql.Where("status = ?", status)
	}
	if types != 0 {
		ql = ql.Where("types = ? ", types)
	}

	if keyword != "" {
		s := "%"
		keyword = s + keyword + s
		ql = ql.Where("title like ? or creator_name like ? ", keyword, keyword)
	}
	var total int64
	ql.Count(&total)
	ql = ql.Limit(limit).Offset((page - 1) * limit)
	ql = ql.Order("created_at DESC")
	messageList := make([]*models.MessageList, 0)
	err := ql.Find(&messageList).Error
	return messageList, total, err
}
