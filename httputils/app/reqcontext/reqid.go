package reqcontext

import (
	"fmt"
	"net/http"

	"github.com/prasannavl/go-gluons/httputils/responder"

	"github.com/google/uuid"
	"github.com/prasannavl/goerror/httperror"
)

const RequestIDHeaderKey = "X-Request-Id"

func CreateRequestIDHandler(reuseUpstreamID bool) middleware {
	if reuseUpstreamID {
		return RequestIDInitOrReuseHandler
	}
	return RequestIDInitOrFailHandler
}

func RequestIDInitOrFailHandler(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		err := requestIDInitOrFailHandler(r, FromRequest(r))
		if err != nil {
			responder.SendErrorText(w, err)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func RequestIDInitOrReuseHandler(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		err := requestIDInitOrReuseHandler(r, FromRequest(r))
		if err != nil {
			responder.SendErrorText(w, err)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func requestIDInitOrFailHandler(r *http.Request, c *RequestContext) error {
	if _, ok := r.Header[RequestIDHeaderKey]; ok {
		msg := fmt.Sprintf("illegal header (%s)", RequestIDHeaderKey)
		return httperror.New(400, msg, true)
	}
	c.RequestID, _ = uuid.NewRandom()
	c.Logger = c.Logger.With("reqid", c.RequestID)
	return nil
}

func requestIDInitOrReuseHandler(r *http.Request, c *RequestContext) error {
	var uid uuid.UUID
	if ok, err := requestIDParseFromHeaders(&uid, r.Header); err != nil {
		msg := fmt.Sprintf("malformed header (%s)", RequestIDHeaderKey)
		return httperror.NewWithCause(400, msg, err, true)
	} else if !ok {
		uid, _ = uuid.NewRandom()
	}
	c.RequestID = uid
	c.Logger = c.Logger.With("reqid", c.RequestID)
	return nil
}

func requestIDParseFromHeaders(target *uuid.UUID, h http.Header) (ok bool, err error) {
	if idStr, ok := h[RequestIDHeaderKey]; ok {
		if *target, err = uuid.Parse(idStr[0]); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
