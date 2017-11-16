package middleware

import (
	"fmt"
	"net/http"

	"github.com/prasannavl/go-gluons/http/reqcontext"
	"github.com/prasannavl/mchain"

	"github.com/google/uuid"
	"github.com/prasannavl/go-errors/httperror"
)

const RequestIDHeaderKey = "X-Request-Id"

func RequestIDMiddleware(reuseUpstreamID bool) mchain.Middleware {
	if reuseUpstreamID {
		return requestIDInitOrReuseMiddleware
	}
	return requestIDInitOrFailMiddleware
}

func requestIDInitOrFailMiddleware(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		err := requestIDInitOrFailHandler(r, reqcontext.FromRequest(r))
		if err != nil {
			return err
		}
		return next.ServeHTTP(w, r)
	}
	return mchain.HandlerFunc(f)
}

func requestIDInitOrReuseMiddleware(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		err := requestIDInitOrReuseHandler(r, reqcontext.FromRequest(r))
		if err != nil {
			return err
		}
		return next.ServeHTTP(w, r)
	}
	return mchain.HandlerFunc(f)
}

func requestIDInitOrFailHandler(r *http.Request, c *reqcontext.RequestContext) error {
	if _, ok := r.Header[RequestIDHeaderKey]; ok {
		msg := fmt.Sprintf("illegal header (%s)", RequestIDHeaderKey)
		return httperror.New(400, msg, true)
	}
	c.RequestID, _ = uuid.NewRandom()
	c.Logger = *c.Logger.With("reqid", c.RequestID)
	return nil
}

func requestIDInitOrReuseHandler(r *http.Request, c *reqcontext.RequestContext) error {
	var uid uuid.UUID
	if ok, err := requestIDParseFromHeaders(&uid, r.Header); err != nil {
		msg := fmt.Sprintf("malformed header (%s)", RequestIDHeaderKey)
		return httperror.NewWithCause(400, msg, err, true)
	} else if !ok {
		uid, _ = uuid.NewRandom()
	}
	c.RequestID = uid
	c.Logger = *c.Logger.With("reqid", c.RequestID)
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
