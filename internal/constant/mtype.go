package constant

// RecordStatus 消息读取状态
type RecordStatus int

const (
	// NotRead 未读
	NotRead RecordStatus = 1
	// AlreadyRead 已读
	AlreadyRead RecordStatus = 2
)

// MessageTypes 系统消息 通知消息
type MessageTypes int

const (
	// SystemSort 系统消息
	SystemSort MessageTypes = 1
	// NoticeSort 通知通告
	NoticeSort MessageTypes = 2
)

// MessageStatus 状态
type MessageStatus int

const (
	// Draft 草稿
	Draft MessageStatus = 1
	// Sending 发送中
	Sending MessageStatus = 2
	// AlreadySent 已发送
	AlreadySent MessageStatus = 3
)

// Receiver receiver 定义
type Receiver struct {
	// Type 1: 人员 2:部门
	Type int    `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
