package client

import (
	"context"
	"io/ioutil"
	"net/http"
)

// FileServerAPI FileServerAPI
type FileServerAPI interface {
	GetFile(ctx context.Context, path string) ([]byte, error)
}

type fileServerAPI struct {
	host string
}

// NewFileServer NewFileServer
func NewFileServer() FileServerAPI {
	return &fileServerAPI{}
}

func (file *fileServerAPI) GetFile(ctx context.Context, path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
