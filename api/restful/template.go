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

// Template Controller
type Template struct {
	template service.Template
	log      logr.Logger
}

// NewTemplate 创建
func NewTemplate(conf *config.Config, log logr.Logger) (*Template, error) {
	log = log.WithName("restful-template")
	t, err := service.NewTemplate(conf, log)
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
	ctx := header.MutateContext(c)

	req := &service.CreateTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.CreateTemplate(ctx, req)).Context(c)
}

// UpdateTemplate update
func (t *Template) UpdateTemplate(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.UpdateTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.UpdateTemplate(ctx, req)).Context(c)
}

// DeleteTemplate 删除
func (t *Template) DeleteTemplate(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.DeleteTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.DeleteTemplate(ctx, req)).Context(c)

}

//QueryTemplatePage 条件查询
func (t *Template) QueryTemplatePage(c *gin.Context) {
	ctx := header.MutateContext(c)

	req := &service.QueryTemplateReq{}
	if err := c.ShouldBind(req); err != nil {
		t.log.Error(err, "should bind", header.GetRequestIDKV(ctx).Fuzzy())
		resp.Format(nil, err).Context(c, http.StatusBadRequest)
		return
	}
	resp.Format(t.template.QueryTemplate(ctx, req)).Context(c)
}
