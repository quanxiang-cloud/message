package code

import "git.internal.yunify.com/qxp/misc/error2"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// ErrDeleteMsState 删除错误的消息
	ErrDeleteMsState = 40014000001
	//ErrNotExistTemplate 模板不存在
	ErrNotExistTemplate = 40014000002
)

// CodeTable 码表
var CodeTable = map[int64]string{
	ErrDeleteMsState:    "该条消息不处于草稿阶段，不能删除",
	ErrNotExistTemplate: "消息模板不存在",
}
