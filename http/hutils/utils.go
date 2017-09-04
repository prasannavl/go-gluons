package hutils

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/prasannavl/mchain"
)

func RunOnPrefix(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		return true, StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func RunOnPrefixAndRedirectToSlash(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		if r.URL.Path == prefix {
			http.RedirectHandler(prefix+"/", http.StatusMovedPermanently).ServeHTTP(w, r)
			return true, nil
		}
		return true, StripPrefix(prefix, h).ServeHTTP(w, r)
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

func StripPrefix(prefix string, h mchain.Handler) mchain.Handler {
	if prefix == "" {
		return h
	}
	return mchain.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		p := strings.TrimPrefix(r.URL.Path, prefix)
		r2 := new(http.Request)
		*r2 = *r
		r2.URL = new(url.URL)
		*r2.URL = *r.URL
		r2.URL.Path = p
		return h.ServeHTTP(w, r2)
	})
}
