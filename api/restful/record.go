package restful

import (
	"github.com/go-logr/logr"
	"net/http"

	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"github.com/gin-gonic/gin"

	"github.com/quanxiang-cloud/message/internal/service"
	"github.com/quanxiang-cloud/message/pkg/config"
)

// Record 消息发送结构体
type Record struct {
	record service.Record
	log    logr.Logger
}

// NewRecord NewRecord
func NewRecord(conf *config.Config, log logr.Logger) (*Record, error) {

	m, err := service.NewRecord(conf)

	if err != nil {
		return nil, err
	}
	return &Record{
		record: m,
		log:    log.WithName("controller record "),
	}, nil
}

// CenterMsByID find message by id
func (m *Record) CenterMsByID(c *gin.Context) {
	req := &service.CenterMsByIDReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.CenterMsByID(logger.CTXTransfer(c), req)).Context(c)
}

//GetNumber  dep  reciver get not read number
func (m *Record) GetNumber(c *gin.Context) {
	req := &service.GetNumberReq{}
	req.ReceiverID = header2.GetProfile(c).UserID
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.GetNumber(logger.CTXTransfer(c), req)).Context(c)

}

// AllRead update already read basis of receiverID
func (m *Record) AllRead(c *gin.Context) {
	req := &service.AllReadReq{}
	req.ReceiverID = header2.GetProfile(c).UserID
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.AllRead(logger.CTXTransfer(c), req)).Context(c)
}

//DeleteByIDs  delete message by IDs
func (m *Record) DeleteByIDs(c *gin.Context) {
	req := &service.DeleteByIDsReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.DeleteByIDs(logger.CTXTransfer(c), req)).Context(c)
}

// ReadByIDs read message by Ids
func (m *Record) ReadByIDs(c *gin.Context) {
	req := &service.ReadByIDsReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.ReadByIDs(logger.CTXTransfer(c), req)).Context(c)
}

// GetMesSendList get by condition
func (m *Record) GetMesSendList(c *gin.Context) {
	req := &service.RecordListReq{}
	req.ReceiverID = header2.GetProfile(c).UserID

	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.RecordList(logger.CTXTransfer(c), req)).Context(c)
}
