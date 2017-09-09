package utils

import (
	"net/http"
	"strings"

	"github.com/prasannavl/go-gluons/http/middleware"
	"github.com/prasannavl/mchain"
)

func OnPrefix(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		return true, h.ServeHTTP(w, r)
	}
	return false, nil
}

func OnPrefixFunc(prefix string, f mchain.HandlerFunc, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	return OnPrefix(prefix, mchain.HandlerFunc(f), w, r)
}

func OnStrippedPrefix(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		return true, middleware.StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func OnStrippedPrefixFunc(prefix string, h mchain.HandlerFunc, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	return OnStrippedPrefix(prefix, mchain.HandlerFunc(h), w, r)
}

func OnStrippedPrefixAndRedirectToSlash(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		if r.URL.Path == prefix {
			// note: Host from the url, which often is empty, unless proxied.
			// This is not the Host header.
			path := r.URL.Host + middleware.ConstructPathFromStripped(r) + "/"
			if r.URL.RawQuery != "" {
				path += "?" + r.URL.RawQuery
			}
			RedirectHandler(
				path,
				http.StatusMovedPermanently).ServeHTTP(w, r)
			return true, nil
		}
		return true, middleware.StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func OnStrippedPrefixAndRedirectToSlashFunc(prefix string, f mchain.HandlerFunc, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	return OnStrippedPrefixAndRedirectToSlash(prefix, mchain.HandlerFunc(f), w, r)
}

// Handlers

func Mount(prefix string, h mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := OnStrippedPrefix(prefix, h, w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func MountAndRedirectToSlash(prefix string, h mchain.Handler) mchain.Handler {
	prefix = strings.TrimSuffix(prefix, "/")
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := OnStrippedPrefixAndRedirectToSlash(prefix, h, w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func MountFunc(prefix string, f mchain.HandlerFunc) mchain.Handler {
	return Mount(prefix, mchain.HandlerFunc(f))
}

func MountFuncAndRedirectToSlash(prefix string, f mchain.HandlerFunc) mchain.Handler {
	return MountAndRedirectToSlash(prefix, mchain.HandlerFunc(f))
}

func HandlerOnPrefix(prefix string, h mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := OnPrefix(prefix, h, w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func HandlerFuncOnPrefix(prefix string, f mchain.HandlerFunc) mchain.Handler {
	return HandlerOnPrefix(prefix, mchain.HandlerFunc(f))
}
