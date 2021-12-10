package service

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mysql2"
	"github.com/quanxiang-cloud/message/internal/core"
	"github.com/quanxiang-cloud/message/internal/models"
	"github.com/quanxiang-cloud/message/internal/models/mysql"
	"github.com/quanxiang-cloud/message/pkg/config"
	"gorm.io/gorm"
)

// Send message send
type Send interface {
	CenterMsByID(ctx context.Context, req *CenterMsByIDReq) (*CenterMsByIDResp, error)
	GetNumber(ctx context.Context, req *GetNumberReq) (*GetNumberResp, error)
	AllRead(ctx context.Context, req *AllReadReq) (*AllReadResp, error)
	DeleteByIDs(ctx context.Context, req *DeleteByIDsReq) (*DeleteByIDsResp, error)
	ReadByIDs(ctx context.Context, req *ReadByIDsReq) (*ReadByIdsResp, error)
	GetMesSendList(ctx context.Context, req *GetMesSendListReq) (*GetMesSendListResp, error)
}

// NewMessageSend create
func NewMessageSend(conf *config.Config) (Send, error) {
	db, err := mysql2.New(conf.Mysql, logger.Logger)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &messageSend{
		conf:            conf,
		db:              db,
		messageSendRepo: mysql.NewMessageSendRepo(),
	}, nil
}

type messageSend struct {
	conf            *config.Config
	db              *gorm.DB
	messageSendRepo models.MessageSendRepo
}

// CenterMsByIDReq req
type CenterMsByIDReq struct {
	ID   string `json:"id"`
	Read bool   `json:"read"`
}

// CenterMsByIDResp resp
type CenterMsByIDResp struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Title   string `json:"title"`

	HandleName string              `json:"handle_name"`
	ReadStatus models.MSReadStatus `json:"read_status"`

	Sort models.MesSort `json:"sort"`

	CreatedAt int64 `json:"created_at"`

	UpdatedAt     int64               `json:"update_at"`
	MesAttachment []models.Attachment `json:"mes_attachment"`
}

// CenterMsByID byID
func (ms *messageSend) CenterMsByID(ctx context.Context, req *CenterMsByIDReq) (*CenterMsByIDResp, error) {
	msSend, err := ms.messageSendRepo.GetByID(ms.db, req.ID)
	if err != nil {
		return nil, err
	}
	// 是否在查询消息的时候，把消息标记为已读
	if req.Read {
		err = ms.messageSendRepo.ReadByID(ms.db, req.ID)
		if err != nil {
			return nil, err
		}
	}
	var resp *CenterMsByIDResp
	if msSend != nil {
		resp = &CenterMsByIDResp{
			ID: msSend.ID,

			Title:         msSend.Title,
			HandleName:    msSend.HandleName,
			CreatedAt:     msSend.CreatedAt,
			UpdatedAt:     msSend.UpdatedAt,
			ReadStatus:    msSend.ReadStatus,
			MesAttachment: msSend.MesAttachment,
			Sort:          msSend.Sort,
		}
		if msSend.Content != "" {
			dataByte, err := base64.StdEncoding.DecodeString(msSend.Content) //  解码
			if err != nil {
				return nil, err
			}
			var ct string
			json.Unmarshal(dataByte, &ct)
			resp.Content = ct
		}
	}
	return resp, nil
}

// GetNumberReq req
type GetNumberReq struct {
	ReciverID string `json:"reciver_id"`
}

// GetNumberResp resp
type GetNumberResp struct {
	TypeNum []*GetNumberRespVO `json:"type_num"`
}

// GetNumberRespVO vo
type GetNumberRespVO struct {
	Total int64          `json:"total"`
	Sort  models.MesSort `json:"sort"`
}

// GetNumber 获取不同消息类型，未读的条数
func (ms *messageSend) GetNumber(ctx context.Context, req *GetNumberReq) (*GetNumberResp, error) {
	numResult, err := ms.messageSendRepo.GetNumber(ms.db, req.ReciverID)
	if err != nil {
		return nil, err
	}
	resp := &GetNumberResp{
		TypeNum: make([]*GetNumberRespVO, len(numResult)),
	}

	for i, rs := range numResult {
		vo := &GetNumberRespVO{
			Total: rs.Total,
			Sort:  rs.Sort,
		}
		resp.TypeNum[i] = vo
	}
	return resp, nil
}

// AllReadReq req
type AllReadReq struct {
	ReciverID string
}

// AllReadResp resp
type AllReadResp struct {
}

// AllRead allread
func (ms *messageSend) AllRead(ctx context.Context, req *AllReadReq) (*AllReadResp, error) {
	err := ms.messageSendRepo.UpdateReadStatus(ms.db, req.ReciverID)
	if err != nil {
		return nil, err
	}
	return &AllReadResp{}, nil
}

// DeleteByIDsReq req
type DeleteByIDsReq struct {
	ArrID []string `json:"arr_id"`
}

// DeleteByIDsResp resp
type DeleteByIDsResp struct {
}

// DeleteByIDs byid
func (ms *messageSend) DeleteByIDs(ctx context.Context, req *DeleteByIDsReq) (*DeleteByIDsResp, error) {
	err := ms.messageSendRepo.DeleteByIDs(ms.db, req.ArrID)
	if err != nil {
		return nil, err
	}
	return &DeleteByIDsResp{}, nil
}

// ReadByIDsReq req
type ReadByIDsReq struct {
	ArrID []string `json:"arr_id"`
}

// ReadByIdsResp resp
type ReadByIdsResp struct {
}

// ReadByIDs readbyids
func (ms *messageSend) ReadByIDs(ctx context.Context, req *ReadByIDsReq) (*ReadByIdsResp, error) {
	err := ms.messageSendRepo.ReadByIDs(ms.db, req.ArrID)
	if err != nil {
		return nil, err
	}
	return &ReadByIdsResp{}, nil
}

// GetMesSendListReq req
type GetMesSendListReq struct {
	ReadStatus int8 `json:"read_status"`
	MesSort    int8 `json:"sort"`
	ReciverID  string
	Channel    string `json:"channel"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// GetMesSendListResp resp
type GetMesSendListResp struct {
	List  []*MesListVO `json:"mes_list"`
	Total int64        `json:"total"`
}

// MesListVO vo
type MesListVO struct {
	ID         string              `json:"id"`
	Title      string              `json:"title"`
	CreatedAt  int64               `json:"updated_at"`
	Sort       models.MesSort      `json:"sort"`
	ReadStatus models.MSReadStatus `json:"read_status"`
}

// GetMesSendList list
func (ms *messageSend) GetMesSendList(ctx context.Context, req *GetMesSendListReq) (*GetMesSendListResp, error) {
	if req.Channel == "" {
		req.Channel = core.Letter.String()
	}
	msList, total, err := ms.messageSendRepo.List(ms.db, req.ReadStatus, req.MesSort, req.Page, req.Limit, req.ReciverID, req.Channel)
	if err != nil {
		return nil, err
	}
	resp := &GetMesSendListResp{
		List: make([]*MesListVO, len(msList)),
	}
	for i, msSend := range msList {
		resp.List[i] = new(MesListVO)
		cloneMsSend(resp.List[i], msSend)
	}
	resp.Total = total
	return resp, nil
}

func cloneMsSend(dst *MesListVO, src *models.MessageSend) {
	dst.ID = src.ID
	dst.Title = src.Title
	dst.CreatedAt = src.CreatedAt
	dst.Sort = src.Sort
	dst.ReadStatus = src.ReadStatus
}
