package util

import (
	"github.com/go-chi/render"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// ErrorResponse defines an error response returned to user
type ErrorResponse struct {
	HTTPStatusCode int    `json:"-"`
	Err            string `json:"error"`
}

// Render fills out a status code defined in ErrorResponse
func (e *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func NewErrorResponse(err error) render.Renderer {

	grpcStatus, ok := status.FromError(err)
	var statusCode int
	if ok {
		code := grpcStatus.Code()
		switch code {
		case codes.Canceled:
			statusCode = http.StatusInternalServerError
		case codes.InvalidArgument:
			statusCode = http.StatusBadRequest
		case codes.NotFound:
			statusCode = http.StatusNotFound
		case codes.AlreadyExists:
			statusCode = http.StatusConflict
		case codes.PermissionDenied:
			statusCode = http.StatusForbidden
		case codes.Unavailable:
			statusCode = http.StatusServiceUnavailable
		case codes.Unauthenticated:
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		return &ErrorResponse{HTTPStatusCode: statusCode, Err: grpcStatus.Message()}
	}
	return &ErrorResponse{HTTPStatusCode: http.StatusInternalServerError, Err: err.Error()}
}
