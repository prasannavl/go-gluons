package chainutils

import (
	"net/http"

	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/hutils"
)

func Mount(prefix string, h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			done, err := hutils.RunOnPrefix(prefix, h, w, r)
			if done {
				return err
			} else {
				return next.ServeHTTP(w, r)
			}
		}
		return mchain.HandlerFunc(f)
	}
	return hh
}

func MountRedirectToSlashed(prefix string, h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			done, err := hutils.RunOnPrefixAndRedirectToSlash(prefix, h, w, r)
			if done {
				return err
			} else {
				return next.ServeHTTP(w, r)
			}
		}
		return mchain.HandlerFunc(f)
	}
	return hh
}

func MountFuncRedirectToSlashed(prefix string, h mchain.HandlerFunc) mchain.Middleware {
	return MountRedirectToSlashed(prefix, mchain.HandlerFunc(h))
}

func MountFunc(prefix string, h mchain.HandlerFunc) mchain.Middleware {
	return Mount(prefix, mchain.HandlerFunc(h))
}
