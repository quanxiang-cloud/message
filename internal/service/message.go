package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/quanxiang-cloud/message/internal/constant"
	"net/http"
	template2 "text/template"

	"git.internal.yunify.com/qxp/misc/header2"
	"github.com/quanxiang-cloud/message/pkg/component/event"
	"gorm.io/gorm"

	"github.com/quanxiang-cloud/message/pkg/client"

	"git.internal.yunify.com/qxp/misc/error2"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mysql2"
	"git.internal.yunify.com/qxp/misc/time2"

	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/internal/models"
	"github.com/quanxiang-cloud/message/internal/models/mysql"
	"github.com/quanxiang-cloud/message/pkg/code"
	"github.com/quanxiang-cloud/message/pkg/config"
)

const (
	send = "http://message/api/v1/message/send"
)

// Message message
type Message interface {
	CreateMessage(ctx context.Context, req *CreateMessageReq) (*CreateMessageResp, error)

	GetMesByID(ctx context.Context, req *GetMesByIDReq) (*GetMesByIDResp, error)

	DeleteMessage(ctx context.Context, req *DeleteMessageReq) (*DeleteMessageResp, error)

	MessageList(ctx context.Context, req *ListReq) (*ListResp, error)
}

type message struct {
	conf         *config.Config
	db           *gorm.DB
	messageRepo  models.MessageRepo
	templateRepo models.TemplateRepo
	recordRepo   models.RecordRepo
	log          logr.Logger
	client       http.Client

	userClient client.User
}

// NewMessage create
func NewMessage(conf *config.Config, log logr.Logger) (Message, error) {
	db, err := mysql2.New(conf.Mysql, logger.Logger)
	if err != nil {
		return nil, err
	}

	return &message{
		conf:         conf,
		db:           db,
		messageRepo:  mysql.NewMessageRepo(),
		templateRepo: mysql.NewTemplateRepo(),
		recordRepo:   mysql.NewRecordRepo(),
		userClient:   client.NewUser(conf.InternalNet),
		log:          log.WithName("service message"),
	}, nil
}

type CreateMessageReq struct {
	Profile header2.Profile
	data    `json:",omitempty"`
}

type data struct {
	Letter *letter `json:"letter"`
	Email  *email  `json:"sms"`
	Web    *web    `json:"web"`
}

type letter struct {
	ID      string   `json:"id,omitempty"`
	UUID    []string `json:"uuid,omitempty"`
	Content *content `json:"contents"`
}

type web struct {
	ID        string                `json:"id"`
	Types     constant.MessageTypes `json:"types" binding:"-"`  // 1. 系统消息 2、 通知通告'
	IsSend    bool                  `json:"isSend" binding:"-"` //  1. draft    2.  send
	Title     string                `json:"title"`
	Files     models.Files          `json:"files" binding:"-"`     // 消息附件
	Receivers models.Receivers      `json:"receivers" binding:"-"` // 接收者
	Content   *content              `json:"contents"`
}

type email struct {
	To          []string `json:"To"`
	Title       string   `json:"title"`
	Content     *content `json:"contents"`
	ContentType string   `json:"content_type,omitempty"`
}

type content struct {
	Content     string            `json:"content"`
	TemplateID  string            `json:"templateID"`
	KeyAndValue map[string]string `json:"keyAndValue"`
}

// CreateMessageResp resp
type CreateMessageResp struct {
	ID string `json:"id"`
}

func (m *message) CreateMessage(ctx context.Context, req *CreateMessageReq) (*CreateMessageResp, error) {
	if req.Letter != nil {
		return m.createLetter(ctx, req.Letter)

	}
	if req.Email != nil {
		return m.createEmail(ctx, req.Email)
	}
	if req.Web != nil {
		return m.createWeb(ctx, req.Web, req.Profile)
	}
	return nil, nil
}

func (m *message) createWeb(ctx context.Context, web *web, profile header2.Profile) (*CreateMessageResp, error) {
	// 只有web 需要入库
	tx := m.db.Begin()
	var err error
	if web.ID != "" {
		err = m.messageRepo.Delete(tx, web.ID)
		if err != nil {
			return nil, err
		}
	}
	convertContent, err := m.convertContent(web.Content)
	if err != nil {
		return nil, err

	}
	messages := &models.MessageList{
		ID:          id2.GenID(),
		Title:       web.Title,
		CreatorID:   profile.UserID,
		CreatorName: profile.UserName,
		Types:       web.Types,
		Status:      constant.Draft,
		Receivers:   web.Receivers,
		CreatedAt:   time2.NowUnix(),
		Files:       web.Files,
		Content:     convertContent,
	}
	err = m.messageRepo.Create(tx, messages)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// 不需要发送， 直接return
	tx.Commit()
	if !web.IsSend {

		return &CreateMessageResp{
			ID: messages.ID,
		}, nil
	}
	m.webSend(ctx, web, messages.ID, convertContent)
	return &CreateMessageResp{
		ID: messages.ID,
	}, nil
}

func (m *message) webSend(ctx context.Context, webData *web, messageID, convertContent string) error {
	var failCount, totalCount int
	for _, value := range webData.Receivers {
		if value.Type == models.Department {
			userInfo, err := m.userClient.GetUsersByDEPID(ctx, value.ID, 0, 1, 1000)
			if err != nil {
				continue
			}
			for _, u := range userInfo {
				totalCount = totalCount + 1
				record := &models.Record{
					ID:           id2.GenID(),
					ListID:       messageID,
					Types:        webData.Types,
					ReceiverID:   u.ID,
					ReceiverName: u.UserName,
					CreatedAt:    time2.NowUnix(),
				}
				err = m.recordCreateAndSend(ctx, record, convertContent)
				if err != nil {
					m.log.Error(err, " dep recordCreateAndSend error", "Request-ID", logger.STDRequestID(ctx))
					failCount = failCount + 1
				}
			}
		} else {
			totalCount = totalCount + 1
			record := &models.Record{
				ID:           id2.GenID(),
				Types:        webData.Types,
				ListID:       messageID,
				ReceiverID:   value.ID,
				ReceiverName: value.Name,
				CreatedAt:    time2.NowUnix(),
			}
			err := m.recordCreateAndSend(ctx, record, convertContent)
			if err != nil {
				m.log.Error(err, "user recordCreateAndSend error", "Request-ID", logger.STDRequestID(ctx))
				failCount = failCount + 1
			}

		}
	}
	update := &models.MessageList{
		ID:      messageID,
		SendNum: totalCount,
		Fail:    failCount,
		Success: totalCount - failCount,
		Status:  constant.AlreadySent,
	}
	err := m.messageRepo.UpdateCount(m.db, update)
	if err != nil {
		m.log.Error(err, "update message count error", "Request-ID", logger.STDRequestID(ctx))
	}
	return nil
}

func (m *message) recordCreateAndSend(ctx context.Context, record *models.Record, content string) error {
	err := m.recordRepo.Create(m.db, record)
	if err != nil {
		return err
	}
	contentByte, _ := json.Marshal(content)
	message := new(event.Data)
	message.LetterSpec = &event.LetterSpec{
		ID:      record.ReceiverID,
		Content: contentByte,
	}
	return m.Send(ctx, message)
}

func (m *message) Send(ctx context.Context, message *event.Data) error {

	return client.POST(ctx, &m.client, "http://message/api/v1/message/send", message, nil)

}

func (m *message) createLetter(ctx context.Context, letter *letter) (*CreateMessageResp, error) {
	convertContent, err := m.convertContent(letter.Content)
	if err != nil {
		return nil, err
	}
	contentByte, _ := json.Marshal(convertContent)
	message := new(event.Data)
	message.LetterSpec = &event.LetterSpec{
		ID:      letter.ID,
		UUID:    letter.UUID,
		Content: contentByte,
	}
	err = m.Send(ctx, message)
	if err != nil {
		return nil, err
	}
	return &CreateMessageResp{}, nil
}

func (m *message) createEmail(ctx context.Context, email *email) (*CreateMessageResp, error) {
	convertContent, err := m.convertContent(email.Content)
	if err != nil {
		return nil, err
	}
	message := new(event.Data)
	message.EmailSpec = &event.EmailSpec{
		To:          email.To,
		Title:       email.Title,
		ContentType: email.ContentType,
		Content:     convertContent,
	}
	err = m.Send(ctx, message)
	if err != nil {
		return nil, err
	}
	return &CreateMessageResp{}, nil
}

func (m *message) convertContent(content *content) (string, error) {
	if content.Content != "" {
		return content.Content, nil
	}
	t, err := m.templateRepo.Get(m.db, content.TemplateID)
	if err != nil {
		return "", err
	}
	t2 := template2.New(t.Content)
	buffer := new(bytes.Buffer)
	err = t2.Execute(buffer, content.KeyAndValue)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

// DeleteMessageReq req
type DeleteMessageReq struct {
	ID string `json:"id"`
}

// DeleteMessageResp resp
type DeleteMessageResp struct {
}

// DeleteMessage delete by id
func (m *message) DeleteMessage(ctx context.Context, req *DeleteMessageReq) (*DeleteMessageResp, error) {
	ms, err := m.messageRepo.Get(m.db, req.ID)
	if err != nil {
		return nil, err
	}
	// 只有在草稿中，才能删除消息
	if ms != nil && ms.Status == constant.Draft {
		err = m.messageRepo.Delete(m.db, req.ID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, error2.NewError(code.ErrDeleteMsState)
	}
	return &DeleteMessageResp{}, nil
}

// GetMesByIDReq by id
type GetMesByIDReq struct {
	ID string `json:"id"`
}

// GetMesByIDResp by id resp
type GetMesByIDResp struct {
	ID        string                `json:"id"`
	Title     string                `json:"title"`
	Types     constant.MessageTypes `json:"types"`
	Receivers models.Receivers      `json:"receivers"`
	Content   string                `json:"content"`
	Files     models.Files          `json:"files"`
	CreatorID string                `json:"creatorId"`
	Success   int                   `json:"success"`
	Fail      int                   `json:"fail"`
	SendNum   int                   `json:"sendNum"`
}

// GetMesByID by id
func (m *message) GetMesByID(ctx context.Context, req *GetMesByIDReq) (resp *GetMesByIDResp, err error) {
	ms, err := m.messageRepo.Get(m.db, req.ID)
	if err != nil {
		return
	}
	resp = &GetMesByIDResp{
		ID:        ms.ID,
		Title:     ms.Title,
		Types:     ms.Types,
		Receivers: ms.Receivers,
		Content:   ms.Content,
		Files:     ms.Files,
		Success:   ms.Success,
		Fail:      ms.Fail,
		SendNum:   ms.SendNum,
	}
	return
}

// ListReq ListReq req
type ListReq struct {
	Status  int8   `json:"status"`
	Sort    int8   `json:"sort"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	KeyWord string `json:"key"`
}

// ListResp ListResp resp
type ListResp struct {
	Messages []*MesVO `json:"messages"`
	Total    int64    `json:"total"`
}

// MesVO vo
type MesVO struct {
	ID          string                `json:"id"`
	Types       constant.MessageTypes `json:"types"`
	Title       string                `json:"title"`
	CreatorName string                `json:"createdName"`
	CreatedAt   int64                 `json:"createdAt"`
	SendNum     int                   `json:"sendNum"`
	Success     int                   `json:"success"`
	Fail        int                   `json:"fail"`
	Files       models.Files          `json:"files"`
}

// MessageList   get message_list by condition
func (m *message) MessageList(ctx context.Context, req *ListReq) (*ListResp, error) {
	ms, total, err := m.messageRepo.List(m.db, req.Status, req.Sort, req.KeyWord, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	resp := &ListResp{
		Messages: make([]*MesVO, len(ms)),
	}
	for i, message := range ms {
		resp.Messages[i] = new(MesVO)
		cloneMs(resp.Messages[i], message)
	}
	resp.Total = total
	return resp, nil
}
func cloneMs(dst *MesVO, src *models.MessageList) {
	dst.ID = src.ID
	dst.Title = src.Title
	dst.CreatedAt = src.CreatedAt
	dst.CreatorName = src.CreatorName
	dst.SendNum = src.SendNum
	dst.Success = src.Success
	dst.Fail = src.Fail
	dst.Files = src.Files
}
