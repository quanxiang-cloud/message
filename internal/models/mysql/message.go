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

// UpdateMessage UpdateMessage
func (m *messageRepo) UpdateMessage(db *gorm.DB, message *models.MessageList) error {
	return db.Table(m.TableName()).Where("id = ?", message.ID).Updates(
		map[string]interface{}{
			"title":    message.Title,
			"args":     message.Args,
			"send_way": message.SendWay,
			"recivers": message.Recivers,
		}).Error

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
func (m *messageRepo) List(db *gorm.DB, status, sort int8, keyword, channel, source string, page, limit int) ([]*models.MessageList, int64, error) {
	ql := db.Table(m.TableName())
	if status != 0 {
		ql = ql.Where("status = ?", status)
	}
	if sort != 0 {
		ql = ql.Where("sort = ? ", sort)
	}
	if channel != "" {
		ql = ql.Where("channel = ? ", channel)
	}
	if keyword != "" {
		s := "%"
		keyword = s + keyword + s
		ql = ql.Where("title like ? or handle_name like ? ", keyword, keyword)
	}
	if source != "" {
		ql = ql.Where("source = ?", source)
	}
	var total int64
	ql.Count(&total)
	ql = ql.Limit(limit).Offset((page - 1) * limit)
	ql = ql.Order("created_at DESC")
	messageList := make([]*models.MessageList, 0)
	err := ql.Find(&messageList).Error
	return messageList, total, err
}

// UpdateCountByID UpdateCountByID
func (m *messageRepo) UpdateCountByID(db *gorm.DB, id string, isSuccess bool) error {
	ql := db.Table(m.TableName()).Where("id = ?", id)
	if isSuccess {
		ql = ql.Update("success", gorm.Expr("`success` + 1"))
	} else {
		ql = ql.Update("fail", gorm.Expr("fail + 1"))
	}
	return ql.Error

}

// UpdateStatus update
func (m *messageRepo) UpdateStatus(db *gorm.DB, id string) error {
	return db.Table(m.TableName()).Where("send_num = success + fail and id = ?", id).Update("status", 3).Error
}

// UpdateSendNum UpdateSendNum
func (m *messageRepo) UpdateSendNum(db *gorm.DB, id string, total int32) error {
	return db.Table(m.TableName()).Where(" id = ?", id).Update("send_num", total).Error

}
