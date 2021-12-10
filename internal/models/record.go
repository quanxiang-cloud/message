package models

import (
	"gorm.io/gorm"
)

// MesSort 分类
type MesSort int8

// MesStatus 状态
type MesStatus int8

// MSReadStatus 消息读取状态
type MSReadStatus int8

const (
	// SystemSort 系统消息
	SystemSort MesSort = 1
	// NoticeSort 通知通告
	NoticeSort MesSort = 2
	// SuccessMes 发送成功
	SuccessMes MesStatus = 1
	// FailMes 发送失败
	FailMes MesStatus = 2
	// NotRead 未读
	NotRead MSReadStatus = 1
	// AlreadyRead 已读
	AlreadyRead MSReadStatus = 2
)

// MessageSend 消息记录
type MessageSend struct {
	ID string

	ListID string

	Content string

	Title string

	HandleID string

	HandleName string

	ReciverID string

	ReciverName string

	ReciverAccount string

	Status MesStatus // 1 fail  2 success

	ReadStatus MSReadStatus // 1 not read    2 read

	Channel string // 发送方式

	Sort MesSort // 1. 系统消息   2、 通知通告

	CreatedAt int64

	UpdatedAt int64

	MesAttachment arrAttachment
}

// Result 未读条数结果集定义
type Result struct {
	Total int64
	Sort  MesSort
}

// MessageSendRepo 消息实体操作接口定义
type MessageSendRepo interface {
	Create(*gorm.DB, *MessageSend) error

	GetByID(*gorm.DB, string) (*MessageSend, error)

	GetNumber(*gorm.DB, string) ([]*Result, error)

	UpdateReadStatus(*gorm.DB, string) error

	DeleteByIDs(*gorm.DB, []string) error

	ReadByIDs(*gorm.DB, []string) error

	List(*gorm.DB, int8, int8, int, int, string, string) ([]*MessageSend, int64, error)

	ReadByID(*gorm.DB, string) error

	UpdateStatus(*gorm.DB, string, MesStatus) error

	GetByCondition(*gorm.DB, string, string) (*MessageSend, error)
}
