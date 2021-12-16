package restful

import (
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"

	"git.internal.yunify.com/qxp/misc/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/internal/service"
	"github.com/quanxiang-cloud/message/pkg/config"
	"net/http"
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

func (m *Message) BatchMessage(c *gin.Context) {
	type batch []*service.CreateMessageReq
	batchData := make(batch, 0)
	err := c.ShouldBindJSON(batchData)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	for _, message := range batchData {
		_, _ = m.message.CreateMessage(c, message)

	}
	resp.Format(struct{}{}, nil).Context(c)

}
