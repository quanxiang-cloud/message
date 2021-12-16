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
		conf:       conf,
		db:         db,
		recordRepo: mysql.NewRecordRepo(),
	}, nil
}

type record struct {
	conf       *config.Config
	db         *gorm.DB
	recordRepo models.RecordRepo
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

	CreatorName string              `json:"creatorName"`
	ReadStatus  constant.ReadStatus `json:"readStatus"`

	CreatedAt int64 `json:"createdAt"`

	UpdatedAt int64        `json:"update_at"`
	Files     models.Files `json:"files"`
}

// CenterMsByID byID
func (ms *record) CenterMsByID(ctx context.Context, req *CenterMsByIDReq) (*CenterMsByIDResp, error) {
	msSend, err := ms.recordRepo.GetByID(ms.db, req.ID)
	if err != nil {
		return nil, err
	}
	// 是否在查询消息的时候，把消息标记为已读
	if req.Read {
		err = ms.recordRepo.ReadByID(ms.db, req.ID)
		if err != nil {
			return nil, err
		}
	}
	var resp *CenterMsByIDResp
	if msSend != nil {
		resp = &CenterMsByIDResp{
			ID:          msSend.ID,
			Title:       "",
			CreatorName: "",
			CreatedAt:   msSend.CreatedAt,

			ReadStatus: msSend.ReadStatus,
			//	MesAttachment: nil ,

		}
	}
	return resp, nil
}

// GetNumberReq req
type GetNumberReq struct {
	ReceiverID string `json:"reciver_id"`
}

// GetNumberResp resp
type GetNumberResp struct {
	TypeNum []*GetNumberRespVO `json:"type_num"`
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
	ReciverID string
}

// AllReadResp resp
type AllReadResp struct {
}

// AllRead allread
func (ms *record) AllRead(ctx context.Context, req *AllReadReq) (*AllReadResp, error) {
	err := ms.recordRepo.UpdateReadStatus(ms.db, req.ReciverID)
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
	ArrID []string `json:"arr_id"`
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
	ReadStatus int8   `json:"readStatus"`
	MesSort    int8   `json:"sort"`
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
	CreatedAt  int64                 `json:"updated_at"`
	Types      constant.MessageTypes `json:"types"`
	ReadStatus constant.ReadStatus   `json:"readStatus"`
}

// RecordList RecordList
func (ms *record) RecordList(ctx context.Context, req *RecordListReq) (*RecordListResp, error) {

	msList, total, err := ms.recordRepo.List(ms.db, req.ReadStatus, req.MesSort, req.Page, req.Limit, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	resp := &RecordListResp{
		List: make([]*RecordVo, len(msList)),
	}
	for i, msSend := range msList {
		resp.List[i] = new(RecordVo)
		cloneMsSend(resp.List[i], msSend)
	}
	resp.Total = total
	return resp, nil
}

func cloneMsSend(dst *RecordVo, src *models.Record) {
	dst.ID = src.ID
	dst.Title = ""
	dst.CreatedAt = src.CreatedAt
	dst.ReadStatus = src.ReadStatus
}
