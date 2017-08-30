package reqcontext

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/prasannavl/go-grab/log"
)

type RequestContext struct {
	RequestID uuid.UUID
	StartTime time.Time
	EndTime   time.Time
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
