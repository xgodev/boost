package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/xgodev/boost/model/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error converts errors to grpc errors
func Error(err error) error {

	if errors.IsNotFound(err) {
		return status.Errorf(codes.NotFound, "%s", err.Error())
	} else if errors.IsNotValid(err) || errors.IsBadRequest(err) {
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	} else if errors.IsServiceUnavailable(err) {
		return status.Errorf(codes.Unavailable, "%s", err.Error())
	} else if errors.IsConflict(err) || errors.IsAlreadyExists(err) {
		return status.Errorf(codes.AlreadyExists, "%s", err.Error())
	} else if errors.IsNotImplemented(err) || errors.IsNotProvisioned(err) {
		return status.Errorf(codes.Unimplemented, "%s", err.Error())
	} else if errors.IsUnauthorized(err) {
		return status.Errorf(codes.Unauthenticated, "%s", err.Error())
	} else if errors.IsForbidden(err) {
		return status.Errorf(codes.PermissionDenied, "%s", err.Error())
	} else {
		switch t := err.(type) {
		case validator.ValidationErrors:
			return status.Errorf(codes.InvalidArgument, "%s", t.Error())
		default:
			return status.Errorf(codes.Internal, "%s", t.Error())
		}
	}
}
