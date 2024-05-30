package http

import (
	"bom-pedido-api/presentation/rest"
	"encoding/json"
	netHttp "net/http"
)

type MuxResponseWriter struct {
	Writer netHttp.ResponseWriter
}

func NewMuxResponseWriter(writer netHttp.ResponseWriter) rest.ResponseWriter {
	return &MuxResponseWriter{Writer: writer}
}

func (writer *MuxResponseWriter) StatusOk(body interface{}) error {
	return writer.Status(200).Body(body)
}

func (writer *MuxResponseWriter) StatusNoContent() rest.ResponseWriter {
	writer.Writer.WriteHeader(204)
	return writer
}

func (writer *MuxResponseWriter) Header(key string, value string) rest.ResponseWriter {
	writer.Writer.Header().Set(key, value)
	return writer
}

func (writer *MuxResponseWriter) Status(status int) rest.ResponseWriter {
	writer.Writer.WriteHeader(status)
	return writer
}

func (writer *MuxResponseWriter) Body(body interface{}) error {
	return json.NewEncoder(writer.Writer).Encode(body)
}

func (writer *MuxResponseWriter) HandleError(err error) error {
	body := map[string]string{"error": err.Error()}
	return writer.Status(500).Body(body)
}
