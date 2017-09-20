package handlerutils

import (
	"net/http"
	"strings"

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

func OnPrefixStripped(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if strings.HasPrefix(r.URL.Path, prefix) {
		return true, StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func OnPrefixStrippedFunc(prefix string, h mchain.HandlerFunc, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	return OnPrefixStripped(prefix, mchain.HandlerFunc(h), w, r)
}

func OnPrefixStrippedAndRedirectToSlash(prefix string, h mchain.Handler, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	// Ensure that the last char of prefix is stripped of "/" before comparison
	ln := len(prefix)
	var px string
	if prefix[ln] == byte('/') {
		px = prefix[:ln-1]
	} else {
		px = prefix
	}
	if strings.HasPrefix(r.URL.Path, px) {
		if r.URL.Path == px {
			RedirectHandler(
				// Use a local redirect, since there's no way to ensure the
				// path isn't modified arbitrarily in order to reconstruct it
				UnsafeRedirectPath(r, "./"),
				http.StatusMovedPermanently).ServeHTTP(w, r)
			return true, nil
		}
		// Ensure that the untrimmed prefix is passed on, so that it can be roundtripped.
		return true, StripPrefix(prefix, h).ServeHTTP(w, r)
	}
	return false, nil
}

func OnPrefixStrippedAndRedirectToSlashFunc(prefix string, f mchain.HandlerFunc, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	return OnPrefixStrippedAndRedirectToSlash(prefix, mchain.HandlerFunc(f), w, r)
}

// Handlers

func Mount(prefix string, h mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := OnPrefixStripped(prefix, h, w, r)
		return err
	}
	return mchain.HandlerFunc(f)
}

func MountAndRedirectToSlash(prefix string, h mchain.Handler) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := OnPrefixStrippedAndRedirectToSlash(prefix, h, w, r)
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
