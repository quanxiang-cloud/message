package models

import (
	"github.com/quanxiang-cloud/message/internal/constant"
	"gorm.io/gorm"
)

// Record 消息记录
type Record struct {
	ID string

	ListID string

	ReceiverID string

	ReceiverName string

	Types constant.MessageTypes

	ReadStatus constant.ReadStatus // 1 not read    2 read

	CreatedAt int64
}

// Result 未读条数结果集定义
type Result struct {
	Total int64
	Types constant.MessageTypes
}

// RecordRepo 消息实体操作接口定义
type RecordRepo interface {
	Create(*gorm.DB, *Record) error

	GetByID(*gorm.DB, string) (*Record, error)

	GetNumber(*gorm.DB, string) ([]*Result, error)

	UpdateReadStatus(*gorm.DB, string) error

	DeleteByIDs(*gorm.DB, []string) error

	ReadByIDs(*gorm.DB, []string) error

	List(*gorm.DB, int8, int8, int, int, string) ([]*Record, int64, error)

	ReadByID(*gorm.DB, string) error

	UpdateStatus(*gorm.DB, string, constant.ReadStatus) error

	GetByCondition(*gorm.DB, string, string) (*Record, error)
}
