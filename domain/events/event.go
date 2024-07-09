package events

type Event struct {
	Id   string      `json:"id"`
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
