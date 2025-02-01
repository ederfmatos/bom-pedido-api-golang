package http

const (
	StatusOK                  Status = 200
	StatusCreated             Status = 201
	StatusNoContent           Status = 204
	StatusBadRequest          Status = 400
	StatusUnAuthorized        Status = 401
	StatusNotFound            Status = 404
	StatusConflict            Status = 409
	StatusUnprocessableEntity Status = 422
	StatusInternalServerError Status = 500
)

type Status int
