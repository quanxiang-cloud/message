package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"git.internal.yunify.com/qxp/misc/redis2"
	"github.com/go-logr/logr"
	"github.com/quanxiang-cloud/message/internal/models"
	"github.com/quanxiang-cloud/message/internal/models/redis"
	"github.com/quanxiang-cloud/message/pkg/client"
	wm "github.com/quanxiang-cloud/message/pkg/component/letter/websocket"
	"github.com/quanxiang-cloud/message/pkg/config"
)

type CacheZone struct {
	ip     string
	cache  models.WSConnetRepo
	log    logr.Logger
	client http.Client
}

func NewCacheZone(ctx context.Context, conf *config.Config, log logr.Logger) (*CacheZone, error) {
	log = log.WithName("cacheZone")

	redisClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		log.Error(err, "new redis client")
		return nil, err
	}

	c := &CacheZone{
		cache:  redis.NewWSConnectRepo(redisClient),
		client: client.New(conf.InternalNet),
		log:    log,
	}

	err = c.setLocalIP()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CacheZone) setLocalIP() error {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		c.log.Error(err, "get net addrs")
		return err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				c.ip = ipnet.IP.String()
				return nil
			}

		}
	}

	return errors.New("can not get local IP")
}

func (c *CacheZone) Create(ctx context.Context, obj wm.Object) error {
	err := c.cache.Create(&models.WSConnet{
		UserID:    obj.ID,
		IP:        c.ip,
		UUID:      obj.UUID,
		CreatedAt: obj.Time.Unix(),
	})
	if err != nil {
		c.log.Error(err, "create", "id", obj.ID, "uuid", obj.UUID)
		return err
	}

	c.cache.Expire(obj.ID)
	return err
}
func (c *CacheZone) Renewal(ctx context.Context, obj wm.Object) error {
	err := c.Create(ctx, obj)
	if err != nil {
		c.log.Error(err, "renewal", "id", obj.ID, "uuid", obj.UUID)
	}
	return err
}

func (c *CacheZone) Delete(ctx context.Context, obj wm.Object) error {
	err := c.cache.Delete(obj.ID, obj.UUID)
	if err != nil {
		c.log.Error(err, "delete", "id", obj.ID, "uuid", obj.UUID)
	}
	return err
}

type PublishReq struct {
	UserID  string   `json:"userID,omitempty"`
	UUID    []string `json:"uuid,omitempty"`
	Content []byte   `json:"content,omitempty"`
}

type PublishResp struct{}

func (c *CacheZone) Publish(ctx context.Context, req *PublishReq) (*PublishResp, error) {
	wsConns, err := c.cache.Get(req.UserID)
	if err != nil {
		c.log.Error(err, "publish", "id", req.UserID, "uuid", req.UUID)
		return &PublishResp{}, err
	}

	boundary := time.Now().Add(-3 * time.Minute).Unix()
	for _, conn := range wsConns {
		// 创建时间距离现在已经超过3分钟了
		if boundary > conn.CreatedAt {
			// 删除 redis 的key
			err := c.cache.Delete(conn.UserID, conn.UUID)
			if err != nil {
				c.log.Error(err, "cache delete", req.UserID, "uuid", req.UUID)
			}
			continue
		}

		if !isSpecific(req.UUID, conn.UUID) {
			continue
		}

		req := wm.SendReq{
			ID:      conn.UserID,
			UUID:    []string{conn.UUID},
			Content: req.Content,
		}

		err = client.POST(ctx, &c.client, fmt.Sprintf("http://%s/api/v1/message/write", conn.IP), req, nil)
		if err != nil {
			// 打印报错信息
			c.log.Error(err, "message send delete fail", req.ID, "uuid", conn.UUID)
		}
	}

	return &PublishResp{}, nil
}

func isSpecific(src []string, specific string) bool {
	if len(src) == 0 {
		return true
	}

	for _, elem := range src {
		if elem == specific {
			return true
		}
	}

	return false
}
