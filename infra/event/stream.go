package event

type Stream interface {
	FetchStream() (chan string, error)
}
