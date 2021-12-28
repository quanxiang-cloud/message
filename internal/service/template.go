package service

import (
	"context"

	"github.com/go-logr/logr"
	id2 "github.com/quanxiang-cloud/cabin/id"
	mysql2 "github.com/quanxiang-cloud/cabin/tailormade/db/mysql"
	"github.com/quanxiang-cloud/cabin/tailormade/header"
	time2 "github.com/quanxiang-cloud/cabin/time"
	"github.com/quanxiang-cloud/message/internal/models"
	"github.com/quanxiang-cloud/message/internal/models/mysql"
	"github.com/quanxiang-cloud/message/pkg/config"
	"gorm.io/gorm"
)

// Template Template
type Template interface {
	CreateTemplate(ctx context.Context, req *CreateTemplateReq) (*CreateTemplateResp, error)

	UpdateTemplate(ctx context.Context, req *UpdateTemplateReq) (*UpdateTemplateResp, error)

	DeleteTemplate(ctx context.Context, req *DeleteTemplateReq) (*DeleteTemplateResp, error)
	QueryTemplate(ctx context.Context, req *QueryTemplateReq) (*QueryTemplateResp, error)
}

// CreateTemplateResp resp
type CreateTemplateResp struct {
}

// UpdateTemplateResp resp
type UpdateTemplateResp struct {
}

// CreateTemplateReq 创建CreateTemplateReq 结构体
type CreateTemplateReq struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateBy string `json:"create_by"`
	Status   int    `json:"status"`
}

// UpdateTemplateReq 修改TemplateReq 结构体
type UpdateTemplateReq struct {
	ID      string `json:"templateId"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  int    `json:"status"`
}

// DeleteTemplateReq 删除
type DeleteTemplateReq struct {
	ID string `json:"id"`
}

// DeleteTemplateResp resp
type DeleteTemplateResp struct {
}

type template struct {
	conf *config.Config
	db   *gorm.DB
	log  logr.Logger

	templateRepo models.TemplateRepo
}

// NewTemplate create
func NewTemplate(conf *config.Config, log logr.Logger) (Template, error) {
	log = log.WithName("service-template")
	db, err := mysql2.New(conf.Mysql, log)
	if err != nil {
		return nil, err
	}
	return &template{
		conf: conf,
		db:   db,
		log:  log,

		templateRepo: mysql.NewTemplateRepo(),
	}, nil
}

// CreateTemplate create
func (t *template) CreateTemplate(ctx context.Context, req *CreateTemplateReq) (*CreateTemplateResp, error) {
	tx := t.db.Begin()
	template := &models.Template{
		ID:        id2.GenID(),
		Name:      req.Name,
		Title:     req.Title,
		Content:   req.Content,
		CreateBy:  req.CreateBy,
		Status:    req.Status,
		CreatedAt: time2.NowUnix(),
		UpdatedAt: time2.NowUnix(),
	}
	err := t.templateRepo.Create(tx, template)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &CreateTemplateResp{}, nil
}

// UpdateTemplate update
func (t *template) UpdateTemplate(ctx context.Context, req *UpdateTemplateReq) (*UpdateTemplateResp, error) {
	template, err := t.templateRepo.Get(t.db, req.ID)
	if err != nil {
		return nil, err
	}
	template.Content = req.Content
	template.UpdatedAt = time2.NowUnix()
	template.Title = req.Title
	template.Name = req.Name
	err = t.templateRepo.UpdateTemplate(t.db, template)
	if err != nil {
		t.log.Error(err, "UpdateTemplate", header.GetRequestIDKV(ctx).Fuzzy()...)
		return nil, err
	}
	return nil, nil
}

// DeleteTemplate delete
func (t *template) DeleteTemplate(ctx context.Context, req *DeleteTemplateReq) (*DeleteTemplateResp, error) {
	err := t.templateRepo.Delete(t.db, req.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// QueryTemplateReq req
type QueryTemplateReq struct {
	Title string `json:"title"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}

// QueryTemplateResp resp
type QueryTemplateResp struct {
	Template []*TemplateVo `json:"template"`
	Total    int64         `josn:"total"`
}

// TemplateVo vo
type TemplateVo struct {
	ID        string `josn:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// QueryTemplate query
func (t *template) QueryTemplate(ctx context.Context, req *QueryTemplateReq) (*QueryTemplateResp, error) {

	templates, total, err := t.templateRepo.QueryTemplate(t.db, req.Title, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	resp := &QueryTemplateResp{

		Template: make([]*TemplateVo, len(templates)),
	}
	for i, value := range templates {
		resp.Template[i] = new(TemplateVo)
		clone(resp.Template[i], value)
	}
	resp.Total = total
	return resp, nil

}
func clone(dst *TemplateVo, src *models.Template) {
	dst.ID = src.ID
	dst.Title = src.Title
	dst.Content = src.Content
	dst.UpdatedAt = src.UpdatedAt
	dst.CreatedAt = src.CreatedAt
}
