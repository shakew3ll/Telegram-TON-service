package apperrors

import (
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Err              error
	StatusCode       codes.Code
	Message          string
	DeveloperMessage string
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}

	return marshal
}

func (e *AppError) GRPCStatus() *status.Status {
	return status.New(e.StatusCode, e.Message)
}

func NewAppError(statusCode codes.Code, message, developerMessage string) *AppError {
	return &AppError{
		StatusCode:       statusCode,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func BadRequestError(message, developerMessage string) *AppError {
	return NewAppError(codes.InvalidArgument, message, developerMessage)
}

func InternalServerError(message, developerMessage string) *AppError {
	return NewAppError(codes.Internal, message, developerMessage)
}

func UnauthorizedError(message, developerMessage string) *AppError {
	return NewAppError(codes.Unauthenticated, message, developerMessage)
}
