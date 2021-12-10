package models

// WSConnet WSConnet
type WSConnet struct {
	UserID    string
	UUID      string
	IP        string
	CreatedAt int64
}

// WSConnetRepo WSConnetRepo
type WSConnetRepo interface {
	Create(*WSConnet) error

	Get(userID string) ([]*WSConnet, error)

	Renewal(userID string) error

	Delete(userID, UUID string) error

	Expire(userID string) error
}
