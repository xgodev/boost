package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/xgodev/boost/model/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error converts a boost/application error to a gRPC status error.
// Errors registered as ignore-as-success return nil.
func Error(err error) error {
	if err == nil {
		return nil
	}
	if pol, ok := errors.IgnoreOf(err); ok && pol&errors.IgnoreAsSuccess != 0 {
		return nil
	}
	if _, ok := err.(validator.ValidationErrors); ok {
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	return status.Errorf(grpcCodeFor(errors.Classify(err)), "%s", err.Error())
}

// grpcCodeFor maps a boost error Kind to a gRPC code.
func grpcCodeFor(kind errors.Kind) codes.Code {
	switch kind {
	case errors.KindNotFound:
		return codes.NotFound
	case errors.KindNotValid, errors.KindBadRequest:
		return codes.InvalidArgument
	case errors.KindServiceUnavailable:
		return codes.Unavailable
	case errors.KindConflict, errors.KindAlreadyExists:
		return codes.AlreadyExists
	case errors.KindNotImplemented, errors.KindNotProvisioned:
		return codes.Unimplemented
	case errors.KindUnauthorized:
		return codes.Unauthenticated
	case errors.KindForbidden:
		return codes.PermissionDenied
	case errors.KindTooManyRequests:
		return codes.ResourceExhausted
	case errors.KindTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Internal
	}
}
