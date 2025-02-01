package handler

import (
	usecase "bom-pedido-api/internal/application/usecase/order"
	"bom-pedido-api/pkg/http"
	"fmt"
)

type OrderHandler struct {
	approveOrderUseCase              *usecase.ApproveOrderUseCase
	cancelOrderUseCase               *usecase.CancelOrderUseCase
	finishOrderUseCase               *usecase.FinishOrderUseCase
	cloneOrderUseCase                *usecase.CloneOrderUseCase
	markOrderDeliveringUseCase       *usecase.MarkOrderDeliveringUseCase
	markOrderInProgressUseCase       *usecase.MarkOrderInProgressUseCase
	markOrderAwaitingWithdrawUseCase *usecase.MarkOrderAwaitingWithdrawUseCase
	markOrderAwaitingDeliveryUseCase *usecase.MarkOrderAwaitingDeliveryUseCase
	rejectOrderUseCase               *usecase.RejectOrderUseCase
}

func NewOrderHandler(
	approveOrderUseCase *usecase.ApproveOrderUseCase,
	cancelOrderUseCase *usecase.CancelOrderUseCase,
	finishOrderUseCase *usecase.FinishOrderUseCase,
	cloneOrderUseCase *usecase.CloneOrderUseCase,
	markOrderDeliveringUseCase *usecase.MarkOrderDeliveringUseCase,
	markOrderInProgressUseCase *usecase.MarkOrderInProgressUseCase,
	markOrderAwaitingWithdrawUseCase *usecase.MarkOrderAwaitingWithdrawUseCase,
	markOrderAwaitingDeliveryUseCase *usecase.MarkOrderAwaitingDeliveryUseCase,
	rejectOrderUseCase *usecase.RejectOrderUseCase,
) *OrderHandler {
	return &OrderHandler{
		approveOrderUseCase:              approveOrderUseCase,
		cancelOrderUseCase:               cancelOrderUseCase,
		finishOrderUseCase:               finishOrderUseCase,
		cloneOrderUseCase:                cloneOrderUseCase,
		markOrderDeliveringUseCase:       markOrderDeliveringUseCase,
		markOrderInProgressUseCase:       markOrderInProgressUseCase,
		markOrderAwaitingWithdrawUseCase: markOrderAwaitingWithdrawUseCase,
		markOrderAwaitingDeliveryUseCase: markOrderAwaitingDeliveryUseCase,
		rejectOrderUseCase:               rejectOrderUseCase,
	}
}

func (h OrderHandler) ApproveOrder(request http.Request, response http.Response) error {
	input := usecase.ApproveOrderUseCaseInput{
		OrderId:    request.PathParam("id"),
		ApprovedBy: request.AuthenticatedUser(),
	}
	if err := h.approveOrderUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("approve order: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) CancelOrder(request http.Request, response http.Response) error {
	input := usecase.CancelOrderInput{
		OrderId:     request.PathParam("id"),
		CancelledBy: request.AuthenticatedUser(),
	}
	if err := h.cancelOrderUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) FinishOrder(request http.Request, response http.Response) error {
	input := usecase.FinishOrderInput{
		OrderId:    request.PathParam("id"),
		FinishedBy: request.AuthenticatedUser(),
	}
	if err := h.finishOrderUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("finish order: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) MarkOrderAwaitingDelivery(request http.Request, response http.Response) error {
	input := usecase.MarkOrderAwaitingDeliveryInput{
		OrderId: request.PathParam("id"),
		By:      request.AuthenticatedUser(),
	}
	if err := h.markOrderAwaitingDeliveryUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("mark order awaiting delivery: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) MarkOrderAwaitingWithdraw(request http.Request, response http.Response) error {
	input := usecase.MarkOrderAwaitingWithdrawInput{
		OrderId: request.PathParam("id"),
		By:      request.AuthenticatedUser(),
	}
	if err := h.markOrderAwaitingWithdrawUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("mark order awaiting withdraw: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) MarkOrderDelivering(request http.Request, response http.Response) error {
	input := usecase.MarkOrderDeliveringInput{
		OrderId: request.PathParam("id"),
		By:      request.AuthenticatedUser(),
	}
	if err := h.markOrderDeliveringUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("mark order delivering: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) MarkOrderInProgress(request http.Request, response http.Response) error {
	input := usecase.MarkOrderInProgressInput{
		OrderId: request.PathParam("id"),
		By:      request.AuthenticatedUser(),
	}
	if err := h.markOrderInProgressUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("mark order in progress: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) RejectOrder(request http.Request, response http.Response) error {
	input := usecase.RejectOrderInput{
		OrderId:    request.PathParam("id"),
		RejectedBy: request.AuthenticatedUser(),
	}
	if err := h.rejectOrderUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("reject order: %w", err)
	}
	return response.NoContent()
}

func (h OrderHandler) CloneOrder(request http.Request, response http.Response) error {
	input := usecase.CloneOrderInput{
		OrderId: request.PathParam("id"),
	}
	if err := h.cloneOrderUseCase.Execute(request.Context(), input); err != nil {
		return fmt.Errorf("clone order: %w", err)
	}
	return response.NoContent()
}
