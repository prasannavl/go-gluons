package handlerutils

import (
	"context"
	"net/http"
	"strings"

	"github.com/prasannavl/mchain"
)

type stripPrefixContextKey struct{}

type StrippedPrefixes = []string

func StrippedPrefixesFromRequest(r *http.Request) StrippedPrefixes {
	v := r.Context().Value(stripPrefixContextKey{})
	if v != nil {
		return v.(StrippedPrefixes)
	}
	return nil
}

func WithStrippedPrefixes(r *http.Request, item StrippedPrefixes) *http.Request {
	c := context.WithValue(r.Context(), stripPrefixContextKey{}, item)
	return r.WithContext(c)
}

func ConstructPathFromStripped(r *http.Request) string {
	prefixes := StrippedPrefixesFromRequest(r)
	var path string
	for _, p := range prefixes {
		path += p
	}
	path += r.URL.Path
	return path
}

func StripPrefix(prefix string, h mchain.Handler) mchain.Handler {
	if prefix == "" {
		return h
	}
	return mchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		// Set context stuff
		s := StrippedPrefixesFromRequest(r)
		r2 := WithStrippedPrefixes(r, append(s, prefix))
		r2.URL.Path = p
		return h.ServeHTTP(w, r2)
	})
}
