package models

// 定义数据库映射结构体
import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

//ReciverType 类型
type ReciverType int8

// MsListStatus 状态
type MsListStatus int8

// MsListType 系统消息 通知消息
type MsListType int8

const (
	// Personnel 人员
	Personnel ReciverType = 1
	// Department 部门
	Department ReciverType = 2
	// DraftStatus 草稿
	DraftStatus MsListStatus = 1
	// SendingStatus 发送中
	SendingStatus MsListStatus = 2
	// AlreadySentStatus 已发送
	AlreadySentStatus MsListStatus = 3
	// VerCode 验证码
	VerCode MsListType = 1
)

// Reciver reciver 定义
type Reciver struct {
	// Type 1: 人员 2:部门
	Type    ReciverType `json:"type,omitempty"`
	ID      string      `json:"id,omitempty"`
	Name    string      `json:"name,omitempty"`
	Account string      `json:"account"`
	Args    []*TeArgs   `json:"args"`
}

// MessageList MessageList
type MessageList struct {
	ID         string
	TemplateID string
	Title      string
	Args       string
	HandleID   string
	HandleName string
	// 站内信、 短信、邮件
	Channel string
	// 1、验证码  2、 非验证码
	Type MsListType
	// 第三方的发送方式
	SendWay string
	//  1、系统消息 2、通知通告
	Sort MesSort
	// 1、草稿  2、发送中  3、 已发送
	Status MsListStatus
	// 接收人
	Recivers string
	// 总人数
	SendNum int64
	// 成功人数
	Success int64
	// 失败人数
	Fail int64
	// 消息附件
	MesAttachment arrAttachment

	Source string

	CreatedAt int64

	UpdatedAt int64
}

type arrAttachment []Attachment

// Attachment Attachment
type Attachment struct {
	Name string `json:"file_name"`
	URL  string `json:"file_url"`
}

// TeArgs mes
type TeArgs struct {
	Key    string      `json:"key"`
	Values interface{} `json:"value" binding:"valueValidator"`
}

// Value 实现方法
func (p arrAttachment) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *arrAttachment) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// MessageRepo MessageRepo
type MessageRepo interface {
	Create(*gorm.DB, *MessageList) error

	UpdateMessage(*gorm.DB, *MessageList) error

	Get(*gorm.DB, string) (*MessageList, error)

	Delete(*gorm.DB, string) error

	List(*gorm.DB, int8, int8, string, string, string, int, int) ([]*MessageList, int64, error)

	UpdateCount(*gorm.DB, *MessageList) error

	UpdateCountByID(*gorm.DB, string, bool) error

	UpdateStatus(*gorm.DB, string) error

	UpdateSendNum(*gorm.DB, string, int32) error
}
