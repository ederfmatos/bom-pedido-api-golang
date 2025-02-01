package net

import (
	"bom-pedido-api/pkg/http"
	"bom-pedido-api/pkg/log"
	"context"
	"fmt"
	log2 "log"
	net "net/http"
	"os"
	"os/signal"
	"syscall"
)

type HTTPServer struct {
	server      *net.Server
	middlewares []http.Middleware
}

func NewHTTPServer(address string) http.Server {
	return &HTTPServer{
		server:      &net.Server{Addr: address},
		middlewares: []http.Middleware{},
	}
}

func (s *HTTPServer) AddMiddleware(middleware http.Middleware) {
	s.middlewares = append(s.middlewares, middleware)
}

func (s *HTTPServer) Get(path string, handler http.Handler, middlewares ...http.Middleware) {
	s.method("GET", path, handler, middlewares)
}

func (s *HTTPServer) Post(path string, handler http.Handler, middlewares ...http.Middleware) {
	s.method("POST", path, handler, middlewares)
}

func (s *HTTPServer) Patch(path string, handler http.Handler, middlewares ...http.Middleware) {
	s.method("PATCH", path, handler, middlewares)
}

func (s *HTTPServer) Put(path string, handler http.Handler, middlewares ...http.Middleware) {
	s.method("PUT", path, handler, middlewares)
}

func (s *HTTPServer) Delete(path string, handler http.Handler, middlewares ...http.Middleware) {
	s.method("DELETE", path, handler, middlewares)
}

func (s *HTTPServer) method(method, path string, handler http.Handler, middlewares []http.Middleware) {
	if middlewares == nil {
		middlewares = make([]http.Middleware, 0, 1)
	}
	net.HandleFunc(fmt.Sprintf("%s /api%s", method, path), func(writer net.ResponseWriter, request *net.Request) {
		var (
			httpRequest  = NewRequest(request)
			httpResponse = NewResponse()
		)

		defer s.writeResponse(writer, httpResponse)

		finalHandler := handler

		for i := len(middlewares) - 1; i >= 0; i-- {
			finalHandler = middlewares[i](finalHandler)
		}
		for i := len(s.middlewares) - 1; i >= 0; i-- {
			finalHandler = s.middlewares[i](finalHandler)
		}

		err := finalHandler(httpRequest, httpResponse)
		if err == nil {
			return
		}

		if httpResponse.Status() != http.StatusOK {
			return
		}

		httpResponse.SetStatus(http.StatusInternalServerError)
		_ = httpResponse.SetBody(http.ErrorResponse{Error: "Ocorreu um erro interno, tente novamente"})
	})
}

func (s *HTTPServer) writeResponse(writer net.ResponseWriter, httpResponse http.Response) {
	writer.WriteHeader(int(httpResponse.Status()))
	headers := httpResponse.Headers()
	for key, value := range headers {
		writer.Header().Set(key, value)
	}
	_, _ = writer.Write(httpResponse.Bytes())
}

func (s *HTTPServer) handleError(response http.Response) {
	if response.Status() != http.StatusOK {
		return
	}

	response.SetStatus(http.StatusInternalServerError)
	_ = response.SetBody(http.ErrorResponse{Error: "Ocorreu um erro interno, tente novamente"})
}

func (s *HTTPServer) Run() {
	err := s.server.ListenAndServe()
	if err != nil {
		log2.Fatal(err)
	}
}

func (s *HTTPServer) AwaitInterruptSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func (s *HTTPServer) Shutdown() {
	log.Info("Shutting down server...")
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Error("Error on shutdown server", err)
	}
}
