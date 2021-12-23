package client

import (
	"context"
	"git.internal.yunify.com/qxp/misc/client"
	"net/http"
)

const (
	fileServer   = "/api/v1/fileserver"
	rangeReadURL = fileServer + "/rangRead"
)

// FileServerConfig FileServerConfig
type FileServerConfig struct {
	InternalNet Config
	Host        string
}

// FileServerAPI FileServerAPI
type FileServerAPI interface {
	RangRead(ctx context.Context, req *RangReadReq) (*RangReadResp, error)
}

type fileServerAPI struct {
	client http.Client
	host   string
}

// NewFileServer NewFileServer
func NewFileServer(conf FileServerConfig) FileServerAPI {
	return &fileServerAPI{
		client: New(conf.InternalNet),
		host:   conf.Host,
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

func (file *fileServerAPI) RangRead(ctx context.Context, req *RangReadReq) (*RangReadResp, error) {
	resp := &RangReadResp{}
	err := client.POST(ctx, &file.client, file.host+rangeReadURL, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
