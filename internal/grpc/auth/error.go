package auth

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NotFoundError(err interface{}) error {
	msg := ""
	switch err.(type) {
	case string:
		msg = err.(string)
	case error:
		msg = err.(error).Error()
	}
	return status.Errorf(codes.NotFound, msg)
}

func AlreadyExistsError(err interface{}) error {
	msg := ""
	switch err.(type) {
	case error:
		msg = err.(error).Error()
	default:
		msg = err.(string)
	}
	return status.Errorf(codes.AlreadyExists, msg)
}

func LogicalError(err interface{}) error {
	msg := ""
	switch err.(type) {
	case error:
		msg = err.(error).Error()
	default:
		msg = err.(string)
	}
	return status.Errorf(codes.FailedPrecondition, msg)
}

func PermissionDeniedError(err interface{}) error {
	msg := ""
	switch err.(type) {
	case error:
		msg = err.(error).Error()
	default:
		msg = err.(string)
	}
	return status.Errorf(codes.PermissionDenied, msg)
}

func ValidationError(key string, err interface{}) error {
	msg := ""
	switch err.(type) {
	case error:
		msg = err.(error).Error()
	default:
		msg = err.(string)
	}
	st, _ := status.New(codes.InvalidArgument, "validation error").WithDetails(&errdetails.BadRequest_FieldViolation{
		Field:       key,
		Description: msg,
	})
	return st.Err()
}
