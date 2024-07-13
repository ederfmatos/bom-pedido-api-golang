package event

type Event struct {
	Id            string      `json:"id"`
	CorrelationId string      `json:"correlationId"`
	Name          string      `json:"name"`
	Data          interface{} `json:"data"`
}
