package reqcontext

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/prasannavl/go-gluons/log"
)

type middleware = func(http.Handler) http.Handler

type RequestContext struct {
	RequestID uuid.UUID
	Logger    log.Logger
}

type requestContextKey struct{}

func FromRequest(r *http.Request) *RequestContext {
	return (r.Context().Value(requestContextKey{})).(*RequestContext)
}

func WithRequestContext(r *http.Request, ctx *RequestContext) *http.Request {
	c := context.WithValue(r.Context(), requestContextKey{}, ctx)
	return r.WithContext(c)
}
