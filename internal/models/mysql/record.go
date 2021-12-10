package mysql

import (
	"github.com/quanxiang-cloud/message/internal/models"
	"gorm.io/gorm"
)

type messageSendRepo struct {
}

// TableName TableName
func (m *messageSendRepo) TableName() string {
	return "message_send"
}

// NewMessageSendRepo NewMessageSendRepo
func NewMessageSendRepo() models.MessageSendRepo {
	return &messageSendRepo{}
}

// Create Create
func (m *messageSendRepo) Create(db *gorm.DB, ms *models.MessageSend) error {
	err := db.Table(m.TableName()).Create(ms).Error
	return err
}

// GetByID GetByID
func (m *messageSendRepo) GetByID(db *gorm.DB, id string) (*models.MessageSend, error) {
	msSend := new(models.MessageSend)
	err := db.Table(m.TableName()).Where("id = ?", id).Find(msSend).Error
	if err != nil {
		return nil, err
	}
	if msSend.ID == "" {
		return nil, nil
	}
	return msSend, nil
}

// GetNumber 获取未读条数的
func (m *messageSendRepo) GetNumber(db *gorm.DB, reciverID string) ([]*models.Result, error) {
	results := make([]*models.Result, 0)

	err := db.Table(m.TableName()).Select(" count(*) as total ,sort ").Where("read_status = 1 and channel = 'letter' and `status` = 1  and reciver_id  = ? ", reciverID).Group("sort").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateReadStatus 把某个人的消息，标记为已读
func (m *messageSendRepo) UpdateReadStatus(db *gorm.DB, reciverID string) error {
	return db.Table(m.TableName()).Where("status = 1  and reciver_id = ?", reciverID).Updates(map[string]interface{}{
		"read_status": 2,
	}).Error
}

// DeleteByIDs DeleteByIDs
func (m *messageSendRepo) DeleteByIDs(db *gorm.DB, arrIds []string) error {
	return db.Table(m.TableName()).Where("id in ?", arrIds).Delete(models.MessageSend{}).Error
}

// ReadByIDs ReadByIDs
func (m *messageSendRepo) ReadByIDs(db *gorm.DB, arrIds []string) error {
	//跟新

	return db.Table(m.TableName()).Where("status = 1  and id in ?", arrIds).Updates(map[string]interface{}{
		"read_status": 2,
	}).Error
}

// List list
func (m *messageSendRepo) List(db *gorm.DB, readStatus int8, sort int8, page, limit int, reciverID, channel string) ([]*models.MessageSend, int64, error) {
	ql := db.Table(m.TableName())
	if readStatus != 0 {
		ql = ql.Where("read_status = ? ", readStatus)
	}
	if sort != 0 {
		ql = ql.Where("sort = ? ", sort)
	}
	if reciverID != "" {
		ql = ql.Where("reciver_id = ? ", reciverID)
	}
	if channel != "" {
		ql = ql.Where("channel = ? ", channel)
	}
	ql = ql.Where("status = 1 ") // 只获取发送成功的消息

	var total int64
	ql.Count(&total)
	ql = ql.Limit(limit).Offset((page - 1) * limit)
	ql = ql.Order("created_at DESC")

	msSendList := make([]*models.MessageSend, 0)
	err := ql.Find(&msSendList).Error
	return msSendList, total, err
}

// ReadByID ReadByID
func (m *messageSendRepo) ReadByID(db *gorm.DB, id string) error {
	return db.Table(m.TableName()).Where("status = 1  and id = ?", id).Updates(map[string]interface{}{
		"read_status": 2,
	}).Error
}

// UpdateStatus UpdateStatus
func (m *messageSendRepo) UpdateStatus(db *gorm.DB, id string, status models.MesStatus) error {
	return db.Table(m.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
	}).Error
}

func (m *messageSendRepo) GetByCondition(db *gorm.DB, listID string, reciverID string) (*models.MessageSend, error) {
	msSend := new(models.MessageSend)
	err := db.Table(m.TableName()).Where("list_id = ? and reciver_id = ? ", listID, reciverID).Find(msSend).Error
	if err != nil {
		return nil, err
	}
	if msSend.ID == "" {
		return nil, nil
	}
	return msSend, nil
}
