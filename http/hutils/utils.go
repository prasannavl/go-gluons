package hutils

import (
	"net/http"
	"strings"

	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/mchain"
)

func RunOnPrefix(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		return true, middleware.StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func RunOnPrefixAndRedirectToSlash(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		if r.URL.Path == prefix {
			path := r.URL.Host + middleware.ConstructPathFromStripped(r) + "/"
			if r.URL.RawQuery != "" {
				path += "?" + r.URL.Query().Encode()
			}
			http.RedirectHandler(
				path,
				http.StatusMovedPermanently).ServeHTTP(w, r)
			return true, nil
		}
		return true, middleware.StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func Mount(prefix string, h mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := RunOnPrefix(prefix, h, w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func MountAndRedirectToSlash(prefix string, h mchain.Handler) mchain.Handler {
	prefix = strings.TrimSuffix(prefix, "/")
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := RunOnPrefixAndRedirectToSlash(prefix, h, w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func MountFunc(prefix string, h mchain.HandlerFunc) mchain.Handler {
	return Mount(prefix, mchain.HandlerFunc(h))
}

func MountFuncAndRedirectToSlash(prefix string, h mchain.HandlerFunc) mchain.Handler {
	return MountAndRedirectToSlash(prefix, mchain.HandlerFunc(h))
}
