package handler

import (
	"bom-pedido-api/internal/application/event"
	"bom-pedido-api/pkg/http"
	"fmt"
)

const (
	_mercadoPago = "MERCADO_PAGO"
	_woovi       = "WOOVI"
)

type (
	CallbackHandler struct {
		eventEmitter event.Emitter
	}

	mercadoPagoCallbackRequest struct {
		Action string `body:"action" json:"action,omitempty"`
	}

	wooviCallbackRequest struct {
		Event  string `body:"event"`
		Charge struct {
			CorrelationID string `body:"correlationID"`
		} `body:"charge"`
	}
)

func NewCallbackHandler(eventEmitter event.Emitter) *CallbackHandler {
	return &CallbackHandler{eventEmitter: eventEmitter}
}

func (h CallbackHandler) MercadoPago(request http.Request, response http.Response) error {
	var callbackRequest mercadoPagoCallbackRequest
	if err := request.Bind(&callbackRequest); err != nil || callbackRequest.Action != "payment.updated" {
		return err
	}

	orderId := request.PathParam("orderId")
	callbackEvent := event.NewPaymentCallbackReceived(_mercadoPago, orderId, callbackRequest.Action)
	if err := h.eventEmitter.Emit(request.Context(), callbackEvent); err != nil {
		return fmt.Errorf("emit event: %w", err)
	}

	return response.NoContent()
}

func (h CallbackHandler) Woovi(request http.Request, response http.Response) error {
	var callbackRequest wooviCallbackRequest
	if err := request.Bind(&callbackRequest); err != nil || callbackRequest.Event != "OPENPIX:CHARGE_COMPLETED" {
		return err
	}

	orderId := request.PathParam("orderId")
	callbackEvent := event.NewPaymentCallbackReceived(_woovi, orderId, "PAYMENT_COMPLETED")
	if err := h.eventEmitter.Emit(request.Context(), callbackEvent); err != nil {
		return fmt.Errorf("emit event: %w", err)
	}

	return response.NoContent()
}
