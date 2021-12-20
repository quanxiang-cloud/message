package mysql

import (
	"github.com/quanxiang-cloud/message/internal/models"
	"gorm.io/gorm"
)

type recordRepo struct {
}

// TableName TableName
func (m *recordRepo) TableName() string {
	return "record"
}

// NewRecordRepo NewRecordRepo
func NewRecordRepo() models.RecordRepo {
	return &recordRepo{}
}

// Create Create
func (m *recordRepo) Create(db *gorm.DB, ms *models.Record) error {
	err := db.Table(m.TableName()).Create(ms).Error
	return err
}

// GetByID GetByID
func (m *recordRepo) GetByID(db *gorm.DB, id string) (*models.Record, error) {
	msSend := new(models.Record)
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
func (m *recordRepo) GetNumber(db *gorm.DB, reciverID string) ([]*models.Result, error) {
	results := make([]*models.Result, 0)

	err := db.Table(m.TableName()).Select(" count(*) as total ,types ").Where("read_status = 1 and receiver_id  = ? ", reciverID).Group("types").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateReadStatus 把某个人的消息，标记为已读
func (m *recordRepo) UpdateReadStatus(db *gorm.DB, receiverID string) error {
	return db.Table(m.TableName()).Where("receiver_id = ?", receiverID).Updates(map[string]interface{}{
		"read_status": 2,
	}).Error
}

// DeleteByIDs DeleteByIDs
func (m *recordRepo) DeleteByIDs(db *gorm.DB, arrIds []string) error {
	return db.Table(m.TableName()).Where("id in ?", arrIds).Delete(models.Record{}).Error
}

// ReadByIDs ReadByIDs
func (m *recordRepo) ReadByIDs(db *gorm.DB, arrIds []string) error {
	//跟新

	return db.Table(m.TableName()).Where("id in ?", arrIds).Updates(map[string]interface{}{
		"read_status": 2,
	}).Error
}

// List list
func (m *recordRepo) List(db *gorm.DB, readStatus int, types int, page, limit int, receiverID string) ([]*models.Record, int64, error) {
	ql := db.Table(m.TableName())
	if readStatus != 0 {
		ql = ql.Where("read_status = ? ", readStatus)
	}
	if types != 0 {
		ql = ql.Where("types = ? ", types)
	}
	if receiverID != "" {
		ql = ql.Where("receiver_id = ? ", receiverID)
	}

	var total int64
	ql.Count(&total)
	ql = ql.Limit(limit).Offset((page - 1) * limit)
	ql = ql.Order("created_at DESC")

	msSendList := make([]*models.Record, 0)
	err := ql.Find(&msSendList).Error
	return msSendList, total, err
}

// ReadByID ReadByID
func (m *recordRepo) ReadByID(db *gorm.DB, id string) error {
	return db.Table(m.TableName()).Where("id = ?", id).Updates(map[string]interface{}{
		"read_status": 2,
	}).Error
}

func (m *recordRepo) GetByCondition(db *gorm.DB, listID string, receiverID string) (*models.Record, error) {
	msSend := new(models.Record)
	err := db.Table(m.TableName()).Where("list_id = ? and receiver_id = ? ", listID, receiverID).Find(msSend).Error
	if err != nil {
		return nil, err
	}
	if msSend.ID == "" {
		return nil, nil
	}
	return msSend, nil
}
