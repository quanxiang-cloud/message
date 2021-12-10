package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/message/internal/models"
)

// NewWSConnectRepo new
func NewWSConnectRepo(client *redis.ClusterClient) models.WSConnetRepo {
	return &wsConnectRepo{
		client: client,
	}
}

const (
	cometTTL = time.Second * 60 * 3
)

type wsConnectRepo struct {
	client *redis.ClusterClient
}

func (*wsConnectRepo) key() string {
	return "message:websocket:comet:"
}

func (w *wsConnectRepo) Create(entity *models.WSConnet) error {
	jsonByte, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	return w.client.HSet(context.Background(), w.key()+entity.UserID, entity.UUID, jsonByte).Err()
}

func (w *wsConnectRepo) Get(userID string) ([]*models.WSConnet, error) {
	result := w.client.HGetAll(context.Background(), w.key()+userID)
	if result.Err() == redis.Nil {
		return nil, nil
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	resp := make([]*models.WSConnet, 0)
	arr, err := result.Result()
	if err != nil {
		return nil, err
	}
	for _, value := range arr {
		var ws *models.WSConnet
		json.Unmarshal([]byte(value), &ws)
		resp = append(resp, ws)
	}

	return resp, nil
}

func (w *wsConnectRepo) Renewal(userID string) error {
	return w.client.
		Expire(context.Background(),
			w.key()+userID,
			cometTTL).
		Err()
}

func (w *wsConnectRepo) Delete(userID, UUID string) error {
	return w.client.HDel(context.Background(), w.key()+userID, UUID).Err()
}

func (w *wsConnectRepo) Expire(userID string) error {
	return w.client.Expire(context.Background(), w.key()+userID, cometTTL).Err()
}
