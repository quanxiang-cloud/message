package restful

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/internal/service"
	"github.com/quanxiang-cloud/message/pkg/config"
	"net/http"
)

// Template Controller
type Template struct {
	template service.Template
	log      logr.Logger
}

// NewTemplate 创建
func NewTemplate(conf *config.Config, log logr.Logger) (*Template, error) {
	t, err := service.NewTemplate(conf)
	if err != nil {
		return nil, err
	}
	return &Template{
		template: t,
		log:      log.WithName("controller template "),
	}, nil
}

// CreateTemplate 生成Template 模板
func (t *Template) CreateTemplate(c *gin.Context) {
	// 去调用
	req := &service.CreateTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.CreateTemplate(logger.CTXTransfer(c), req)).Context(c)
}

// UpdateTemplate update
func (t *Template) UpdateTemplate(c *gin.Context) {
	req := &service.UpdateTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.UpdateTemplate(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteTemplate 删除
func (t *Template) DeleteTemplate(c *gin.Context) {
	req := &service.DeleteTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.DeleteTemplate(logger.CTXTransfer(c), req)).Context(c)

}

//QueryTemplatePage 条件查询
func (t *Template) QueryTemplatePage(c *gin.Context) {
	req := &service.QueryTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", "requestID", logger.GINRequestID(c).String)
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.QueryTemplate(logger.CTXTransfer(c), req)).Context(c)
}
