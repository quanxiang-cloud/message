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

// Record 消息发送结构体
type Record struct {
	record service.Record
	log    logr.Logger
}

// NewRecord NewRecord
func NewRecord(conf *config.Config, log logr.Logger) (*Record, error) {
	log = log.WithName("restful-record")

	m, err := service.NewRecord(conf, log)
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
	ctx := header.MutateContext(c)

	req := &service.CenterMsByIDReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.CenterMsByID(ctx, req)).Context(c)
}

//GetNumber  dep  reciver get not read number
func (m *Record) GetNumber(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.GetNumberReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.ReceiverID = c.GetHeader("User-Id")
	resp.Format(m.record.GetNumber(ctx, req)).Context(c)

}

// AllRead update already read basis of receiverID
func (m *Record) AllRead(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.AllReadReq{}
	req.ReceiverID = c.GetHeader("User-Id")
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.AllRead(ctx, req)).Context(c)
}

//DeleteByIDs  delete message by IDs
func (m *Record) DeleteByIDs(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.DeleteByIDsReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.DeleteByIDs(ctx, req)).Context(c)
}

// ReadByIDs read message by Ids
func (m *Record) ReadByIDs(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.ReadByIDsReq{}
	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.ReadByIDs(ctx, req)).Context(c)
}

// GetMesSendList get by condition
func (m *Record) GetMesSendList(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.RecordListReq{}
	req.ReceiverID = c.GetHeader("User-Id")

	if err := c.ShouldBind(req); err != nil {
		m.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.record.RecordList(ctx, req)).Context(c)
}
