package handler

import (
	"bom-pedido-api/pkg/http"
	"context"
	"errors"
	"sync"
)

type (
	Pingable interface {
		Ping(context.Context) error
	}

	HealthHandler struct {
		dependencies []Pingable
	}

	HealthResponse struct {
		Ok bool `json:"ok"`
	}
)

func NewHealthHandler(dependencies ...Pingable) *HealthHandler {
	return &HealthHandler{dependencies: dependencies}
}

func (h HealthHandler) Health(request http.Request, response http.Response) error {
	var (
		errs      = make([]error, 0)
		mutex     = sync.Mutex{}
		waitGroup = sync.WaitGroup{}
	)

	waitGroup.Add(len(h.dependencies))
	for _, dependency := range h.dependencies {
		go h.health(dependency, &mutex, &waitGroup, &errs, request.Context())
	}
	waitGroup.Wait()

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return response.OK(HealthResponse{Ok: true})
}

func (h HealthHandler) health(pingable Pingable, mutex *sync.Mutex, waitGroup *sync.WaitGroup, errs *[]error, ctx context.Context) {
	defer waitGroup.Done()
	err := pingable.Ping(ctx)
	if err == nil {
		return
	}

	mutex.Lock()
	*errs = append(*errs, err)
	mutex.Unlock()
}
