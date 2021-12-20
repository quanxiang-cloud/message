package models

// 定义数据库映射结构体
import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"

	"github.com/quanxiang-cloud/message/internal/constant"
)

type ReceiverType int64

const Department ReceiverType = 2

// MessageList MessageList
type MessageList struct {
	ID    string
	Title string

	Content string

	CreatorID string

	CreatorName string

	//  1、系统消息 2、通知通告
	Types constant.MessageTypes
	// 1、草稿  2、发送中  3、 已发送
	Status constant.MessageStatus
	// 接收人
	Receivers Receivers
	// 总人数
	SendNum int
	// 成功人数
	Success int
	// 失败人数
	Fail int
	// 消息附件
	Files Files

	CreatedAt int64

	UpdatedAt int64
}

// Files Files
type Files []*File

// Receiver receiver 定义
type Receiver struct {
	// Type 1: 人员 2:部门
	Type ReceiverType `json:"type,omitempty"`
	ID   string       `json:"id,omitempty"`
	Name string       `json:"name,omitempty"`
}

type File struct {
	FileName string `json:"fileName,omitempty"`
	URL      string `json:"url,omitempty"`
}

// Receivers Receivers
type Receivers []*Receiver

// Value 实现方法
func (p Files) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *Files) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// Value 实现方法
func (p Receivers) Value() (driver.Value, error) {
	return json.Marshal(p)
}

// Scan 实现方法
func (p *Receivers) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &p)
}

// MessageRepo MessageRepo
type MessageRepo interface {
	Create(*gorm.DB, *MessageList) error

	Get(*gorm.DB, string) (*MessageList, error)

	Delete(*gorm.DB, string) error

	List(*gorm.DB, int8, int8, string, int, int) ([]*MessageList, int64, error)

	UpdateCount(*gorm.DB, *MessageList) error
}
