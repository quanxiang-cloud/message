package dapr

import "errors"

type DaprEvent struct {
	Topic           string `json:"topic"`
	Pubsubname      string `json:"pubsubname"`
	Traceid         string `json:"traceid"`
	ID              string `json:"id"`
	Datacontenttype string `json:"datacontenttype"`
	Data            Data   `json:"data"`
	Type            string `json:"type"`
	Specversion     string `json:"specversion"`
	Source          string `json:"source"`
}

type Data struct {
	*LetterSpec `json:"letter,omitempty"`
	*EmailSpec  `json:"email,omitempty"`
}

type LetterSpec struct {
	ID      string   `json:"id,omitempty"`
	UUID    []string `json:"uuid,omitempty"`
	Content []byte   `json:"content,omitempty"`
}

type EmailSpec struct{}

var (
	ErrDataIsNil = errors.New("data is nil")
)
