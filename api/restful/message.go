package restful

import (
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"

	"net/http"

	"git.internal.yunify.com/qxp/misc/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
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
	req := &service.CreateMessageReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	req.Profile = header2.GetProfile(c)
	resp.Format(m.message.CreateMessage(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteMessage delete message_list by id
func (m *Message) DeleteMessage(c *gin.Context) {
	req := &service.DeleteMessageReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c))
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(m.message.DeleteMessage(logger.CTXTransfer(c), req)).Context(c)
}

// MessageList  get message_list by condition
func (m *Message) MessageList(c *gin.Context) {
	req := &service.ListReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c))
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(m.message.MessageList(logger.CTXTransfer(c), req)).Context(c)
}

// GetMessageByID 根据id得到返回值
func (m *Message) GetMessageByID(c *gin.Context) {
	req := &service.GetMesByIDReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c))
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(m.message.GetMesByID(logger.CTXTransfer(c), req)).Context(c)
}

// BatchMessage 批量发消息
func (m *Message) BatchMessage(c *gin.Context) {
	var batch []service.CreateMessageReq
	if err := c.ShouldBind(&batch); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c))
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	for _, message := range batch {
		_, err := m.message.CreateMessage(c, &message)
		if err != nil {
			m.log.Error(err, "send message ", "requestID", logger.GINRequestID(c))
		}

	}
	resp.Format(struct{}{}, nil).Context(c)

}
