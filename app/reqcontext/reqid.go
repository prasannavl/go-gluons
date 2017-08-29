package reqcontext

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/prasannavl/goerror/httperror"
	"github.com/prasannavl/mchain"
)

const RequestIDHeaderKey = "X-Request-Id"

func CreateReqIDHandler(reuseUpstreamID bool) mchain.SimpleMiddleware {
	if reuseUpstreamID {
		return mchain.SimpleMiddleware(ReqIDInitOrReuseHandler)
	} else {
		return mchain.SimpleMiddleware(ReqIDInitOrFailHandler)
	}
}

func ReqIDInitOrFailHandler(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
	err := reqIDInitOrFailHandler(r, FromRequest(r))
	if err != nil {
		return err
	}
	return next.ServeHTTP(w, r)
}

func ReqIDInitOrReuseHandler(w http.ResponseWriter, r *http.Request, next mchain.Handler) error {
	err := reqIDInitOrReuseHandler(r, FromRequest(r))
	if err != nil {
		return err
	}
	return next.ServeHTTP(w, r)
}

func reqIDInitOrFailHandler(r *http.Request, c *RequestContext) error {
	if _, ok := r.Header[RequestIDHeaderKey]; ok {
		msg := fmt.Sprintf("illegal header (%s)", RequestIDHeaderKey)
		return httperror.New(400, msg, true)
	}
	c.RequestID, _ = uuid.NewRandom()
	c.Logger = c.Logger.WithContext("reqid", c.RequestID)
	return nil
}

func reqIDInitOrReuseHandler(r *http.Request, c *RequestContext) error {
	var uid uuid.UUID
	if ok, err := reqIDParseFromHeaders(&uid, r.Header); err != nil {
		msg := fmt.Sprintf("malformed header (%s)", RequestIDHeaderKey)
		return httperror.NewWithCause(400, msg, err, true)
	} else if !ok {
		uid, _ = uuid.NewRandom()
	}
	c.RequestID = uid
	c.Logger = c.Logger.WithContext("reqid", c.RequestID)
	return nil
}

func reqIDParseFromHeaders(target *uuid.UUID, h http.Header) (ok bool, err error) {
	if idStr, ok := h[RequestIDHeaderKey]; ok {
		if *target, err = uuid.Parse(idStr[0]); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
