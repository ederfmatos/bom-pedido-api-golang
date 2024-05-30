package http

import (
	"bom-pedido-api/presentation/rest"
	"log/slog"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewHttpServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (server *Server) Run(port string) error {
	return http.ListenAndServe(port, server.mux)
}

func (server *Server) HandleFunc(pattern string, handler rest.RequestHandler) {
	server.mux.HandleFunc(pattern, server.requestWrapper(handler))
}

func (server *Server) requestWrapper(requestHandler rest.RequestHandler) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		httpRequest := NewMuxHttpRequest(request)
		responseWriter := NewMuxResponseWriter(writer)
		err := requestHandler(httpRequest, responseWriter)
		if err != nil {
			slog.Error("Ocorreu um erro na requisição", err)
			responseWriter.Status(500).Body(map[string]string{"error": err.Error()})
		}
	}
}
