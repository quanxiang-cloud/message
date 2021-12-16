package service

import (
	"context"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mysql2"
	"github.com/quanxiang-cloud/message/internal/constant"
	"github.com/quanxiang-cloud/message/internal/models"
	"github.com/quanxiang-cloud/message/internal/models/mysql"
	"github.com/quanxiang-cloud/message/pkg/config"
	"gorm.io/gorm"
)

// Record message send
type Record interface {
	CenterMsByID(ctx context.Context, req *CenterMsByIDReq) (*CenterMsByIDResp, error)
	GetNumber(ctx context.Context, req *GetNumberReq) (*GetNumberResp, error)
	AllRead(ctx context.Context, req *AllReadReq) (*AllReadResp, error)
	DeleteByIDs(ctx context.Context, req *DeleteByIDsReq) (*DeleteByIDsResp, error)
	ReadByIDs(ctx context.Context, req *ReadByIDsReq) (*ReadByIdsResp, error)
	RecordList(ctx context.Context, req *RecordListReq) (*RecordListResp, error)
}

// NewRecord create
func NewRecord(conf *config.Config) (Record, error) {
	db, err := mysql2.New(conf.Mysql, logger.Logger)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &record{
		conf:        conf,
		db:          db,
		recordRepo:  mysql.NewRecordRepo(),
		messageRepo: mysql.NewMessageRepo(),
	}, nil
}

type record struct {
	conf        *config.Config
	db          *gorm.DB
	recordRepo  models.RecordRepo
	messageRepo models.MessageRepo
}

// CenterMsByIDReq req
type CenterMsByIDReq struct {
	ID   string `json:"id"`
	Read bool   `json:"read"`
}

// CenterMsByIDResp resp
type CenterMsByIDResp struct {
	ID          string              `json:"id"`
	Content     string              `json:"content"`
	Title       string              `json:"title"`
	CreatorName string              `json:"creatorName"`
	ReadStatus  constant.ReadStatus `json:"readStatus"`
	UpdatedAt   int64               `json:"updateAt"`
	Files       models.Files        `json:"files"`
}

// CenterMsByID byID
func (ms *record) CenterMsByID(ctx context.Context, req *CenterMsByIDReq) (*CenterMsByIDResp, error) {
	record, err := ms.recordRepo.GetByID(ms.db, req.ID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return &CenterMsByIDResp{}, nil
	}
	// 是否在查询消息的时候，把消息标记为已读
	if req.Read {
		err = ms.recordRepo.ReadByID(ms.db, req.ID)
		if err != nil {
			return nil, err
		}
	}
	message, err := ms.messageRepo.Get(ms.db, record.ListID)
	if err != nil {
		return nil, err
	}
	resp := &CenterMsByIDResp{
		ID:          record.ID,
		Title:       message.Title,
		CreatorName: message.CreatorName,
		UpdatedAt:   record.CreatedAt,
		ReadStatus:  record.ReadStatus,
		Files:       message.Files,
	}
	return resp, nil
}

// GetNumberReq req
type GetNumberReq struct {
	ReceiverID string `json:"receiverId"`
}

// GetNumberResp resp
type GetNumberResp struct {
	TypeNum []*GetNumberRespVO `json:"typeNum"`
}

// GetNumberRespVO vo
type GetNumberRespVO struct {
	Total int64                 `json:"total"`
	Types constant.MessageTypes `json:"types"`
}

// GetNumber 获取不同消息类型，未读的条数
func (ms *record) GetNumber(ctx context.Context, req *GetNumberReq) (*GetNumberResp, error) {
	numResult, err := ms.recordRepo.GetNumber(ms.db, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	resp := &GetNumberResp{
		TypeNum: make([]*GetNumberRespVO, len(numResult)),
	}

	for i, rs := range numResult {
		vo := &GetNumberRespVO{
			Total: rs.Total,
			Types: rs.Types,
		}
		resp.TypeNum[i] = vo
	}
	return resp, nil
}

// AllReadReq req
type AllReadReq struct {
	ReceiverID string
}

// AllReadResp resp
type AllReadResp struct {
}

// AllRead allread
func (ms *record) AllRead(ctx context.Context, req *AllReadReq) (*AllReadResp, error) {
	err := ms.recordRepo.UpdateReadStatus(ms.db, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	return &AllReadResp{}, nil
}

// DeleteByIDsReq req
type DeleteByIDsReq struct {
	ArrID []string `json:"ids"`
}

// DeleteByIDsResp resp
type DeleteByIDsResp struct {
}

// DeleteByIDs byID
func (ms *record) DeleteByIDs(ctx context.Context, req *DeleteByIDsReq) (*DeleteByIDsResp, error) {
	err := ms.recordRepo.DeleteByIDs(ms.db, req.ArrID)
	if err != nil {
		return nil, err
	}
	return &DeleteByIDsResp{}, nil
}

// ReadByIDsReq req
type ReadByIDsReq struct {
	ArrID []string `json:"ids"`
}

// ReadByIdsResp resp
type ReadByIdsResp struct {
}

// ReadByIDs readByIDs
func (ms *record) ReadByIDs(ctx context.Context, req *ReadByIDsReq) (*ReadByIdsResp, error) {
	err := ms.recordRepo.ReadByIDs(ms.db, req.ArrID)
	if err != nil {
		return nil, err
	}
	return &ReadByIdsResp{}, nil
}

// RecordListReq req
type RecordListReq struct {
	ReadStatus int    `json:"readStatus"`
	Types      int    `json:"types"`
	ReceiverID string `json:"receiverID"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

// RecordListResp resp
type RecordListResp struct {
	List  []*RecordVo `json:"list"`
	Total int64       `json:"total"`
}

// RecordVo vo
type RecordVo struct {
	ID         string                `json:"id"`
	Title      string                `json:"title"`
	CreatedAt  int64                 `json:"createdAt"`
	Content    string                `json:"content"`
	Types      constant.MessageTypes `json:"types"`
	ReadStatus constant.ReadStatus   `json:"readStatus"`
	Files      models.Files          `json:"files"`
}

// RecordList RecordList
func (ms *record) RecordList(ctx context.Context, req *RecordListReq) (*RecordListResp, error) {

	msList, total, err := ms.recordRepo.List(ms.db, req.ReadStatus, req.Types, req.Page, req.Limit, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	resp := &RecordListResp{
		List: make([]*RecordVo, len(msList)),
	}

	for i, msSend := range msList {
		resp.List[i] = new(RecordVo)
		mslist, err := ms.messageRepo.Get(ms.db, msSend.ListID)
		if err != nil {
			return nil, err
		}
		cloneMsSend(resp.List[i], msSend, mslist)
	}
	resp.Total = total
	return resp, nil
}

func cloneMsSend(dst *RecordVo, src *models.Record, message *models.MessageList) {
	dst.ID = src.ID
	dst.Title = message.Title
	dst.Types = src.Types
	dst.Content = message.Content
	dst.CreatedAt = src.CreatedAt
	dst.ReadStatus = src.ReadStatus
	dst.Files = message.Files
}
