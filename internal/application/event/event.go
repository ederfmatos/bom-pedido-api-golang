package event

type (
	Name string

	Event struct {
		Id            string            `json:"id"`
		CorrelationId string            `json:"correlationId"`
		Name          Name              `json:"name"`
		Data          map[string]string `json:"data"`
	}
)
