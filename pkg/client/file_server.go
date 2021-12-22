package client

import (
	"context"
	"git.internal.yunify.com/qxp/misc/client"
	"net/http"
)

const (
	fileServer   = "http://fileserver/api/v1/fileserver"
	rangeReadURL = fileServer + "/rangRead"
)

// FileServerAPI FileServerAPI
type FileServerAPI interface {
	RangRead(ctx context.Context, path, opt string, offset, size int) (*RangReadResp, error)
}

type fileServerAPI struct {
	client http.Client
}

// NewFileServer NewFileServer
func NewFileServer(conf Config) FileServerAPI {
	return &fileServerAPI{
		client: New(conf),
	}
}

// RangReadResp RangReadResp
type RangReadResp struct {
	Content []byte `json:"content"`
}

// RangReadReq RangReadReq
type RangReadReq struct {
	Path   string `json:"path"`
	Opt    string `json:"opt"`
	Offset int    `json:"offset"`
	Size   int    `json:"size"`
}

func (file *fileServerAPI) RangRead(ctx context.Context, path, opt string, offset, size int) (*RangReadResp, error) {
	params := &RangReadReq{
		Path:   path,
		Opt:    opt,
		Offset: offset,
		Size:   size,
	}
	resp := &RangReadResp{}
	err := client.POST(ctx, &file.client, rangeReadURL, params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
