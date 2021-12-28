package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/message/internal/service"
	"github.com/quanxiang-cloud/message/pkg/config"
)

// Message 消息
type Message struct {
	message service.Message
	log     logr.Logger
}

// NewMessage createNewMessage
func NewMessage(conf *config.Config, log logr.Logger) (*Message, error) {
	m, err := service.NewMessage(conf, log)
	if err != nil {
		return nil, err
	}
	return &Message{
		message: m,
		log:     log.WithName("controller message "),
	}, nil
}

// SaveMessage 保存消息，但是不发送消息、在上一次草稿的基础之上，直接发送消息，直接发送消息，
func (m *Message) SaveMessage(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.CreateMessageReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy()...)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	req.UserID = c.GetHeader("User-Id")
	req.UserName = c.GetHeader("User-Name")
	resp.Format(m.message.CreateMessage(ctx, req)).Context(c)
}

// DeleteMessage delete message_list by id
func (m *Message) DeleteMessage(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.DeleteMessageReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy()...)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(m.message.DeleteMessage(ctx, req)).Context(c)
}

// MessageList  get message_list by condition
func (m *Message) MessageList(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.ListReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy()...)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(m.message.MessageList(ctx, req)).Context(c)
}

// GetMessageByID 根据id得到返回值
func (m *Message) GetMessageByID(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.GetMesByIDReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy()...)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(m.message.GetMesByID(ctx, req)).Context(c)
}

// BatchMessage 批量发消息
func (m *Message) BatchMessage(c *gin.Context) {
	ctx := header.MutateContext(c)

	var batch []service.CreateMessageReq
	if err := c.ShouldBind(&batch); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy()...)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	for _, message := range batch {
		_, err := m.message.CreateMessage(c, &message)
		if err != nil {
			m.log.Error(err, "send message ", header.GetRequestIDKV(ctx).Fuzzy()...)
		}

	}
	resp.Format(struct{}{}, nil).Context(c)

}
