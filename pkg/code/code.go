package code

import "git.internal.yunify.com/qxp/misc/error2"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// ErrDeleteMsState 删除错误的消息
	ErrDeleteMsState = 40014000001
	//ErrNotExistTemplateMsState 模板不存在
	ErrNotExistTemplateMsState = 40014000002
)

// CodeTable 码表
var CodeTable = map[int]string{
	ErrDeleteMsState:           "该条消息不处于草稿阶段，不能删除",
	ErrNotExistTemplateMsState: "消息模板不存在",
}
