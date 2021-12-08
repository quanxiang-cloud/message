package dapr

type DaprEvent struct {
	Topic           string `json:"topic"`
	Pubsubname      string `json:"pubsubname"`
	Traceid         string `json:"traceid"`
	ID              string `json:"id"`
	Datacontenttype string `json:"datacontenttype"`
	Type            string `json:"type"`
	Data            Data   `json:"data"`
	Specversion     string `json:"specversion"`
	Source          string `json:"source"`
}

type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
