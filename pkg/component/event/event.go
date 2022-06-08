package event

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
	*Multiple   `json:"multiple,omitempty"`
}

type Multiple struct {
	Kind string                 `json:"kind,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type LetterSpec struct {
	ID      string   `json:"id,omitempty"`
	UUID    []string `json:"uuid,omitempty"`
	Content []byte   `json:"content,omitempty"`
}

type Attachment struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

type EmailSpec struct {
	To          []string     `json:"to,omitempty"`
	Title       string       `json:"title,omitempty"`
	ContentType string       `json:"content_type,omitempty"`
	Content     string       `json:"content,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

var (
	ErrDataIsNil = errors.New("data is nil")
)
