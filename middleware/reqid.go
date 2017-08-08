package middleware

import (
	"fmt"
	"net/http"
	"pvl/apicore/reqcontext"

	"github.com/google/uuid"
	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"
)

const RequestIDHeaderKey = "X-Request-Id"

func RequestIDMustInitHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		return requestIDMustInitHandler(w, r, next, reqcontext.FromRequest(r))
	}
	return mchain.HandlerFunc(f)
}

func RequestIDAdoptFromHeaderOrInitHandler(next mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		return requestIDAdoptFromHeaderOrInitHandler(w, r, next, reqcontext.FromRequest(r))
	}
	return mchain.HandlerFunc(f)
}

func requestIDMustInitHandler(w http.ResponseWriter, r *http.Request, next mchain.Handler, c *reqcontext.RequestContext) error {
	if _, ok := r.Header[RequestIDHeaderKey]; ok {
		msg := fmt.Sprintf("illegal header (%s)", RequestIDHeaderKey)
		return httperror.New(400, msg, true)
	}
	c.RequestID = mustNewUUID()
	return next.ServeHTTP(w, r)
}

func requestIDAdoptFromHeaderOrInitHandler(w http.ResponseWriter, r *http.Request, next mchain.Handler, c *reqcontext.RequestContext) error {
	var uid uuid.UUID
	if ok, err := requestIDFromHeader(r, &uid); err != nil {
		msg := fmt.Sprintf("malformed header (%s)", RequestIDHeaderKey)
		return httperror.NewWithCause(400, msg, err, true)
	} else if !ok {
		uid = mustNewUUID()
	}
	c.RequestID = uid
	return next.ServeHTTP(w, r)
}

func requestIDFromHeader(r *http.Request, target *uuid.UUID) (ok bool, err error) {
	if idStr, ok := r.Header[RequestIDHeaderKey]; ok {
		if *target, err = uuid.Parse(idStr[0]); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func mustNewUUID() uuid.UUID {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(&err)
	}
	return id
}
