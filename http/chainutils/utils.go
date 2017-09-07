package chainutils

import (
	"net/http"

	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/http/utils"
)

func Hook(h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			err := h.ServeHTTP(w, r)
			if err != nil {
				return err
			}
			return next.ServeHTTP(w, r)
		}
		return mchain.HandlerFunc(f)
	}
	return hh
}

func Run(h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			return h.ServeHTTP(w, r)
		}
		return mchain.HandlerFunc(f)
	}
	return hh
}

func OnPrefix(prefix string, h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			done, err := utils.OnPrefix(prefix, h, w, r)
			if done {
				return err
			}
			return next.ServeHTTP(w, r)
		}
		return mchain.HandlerFunc(f)
	}
	return hh
}

func OnPrefixFunc(prefix string, h mchain.HandlerFunc) mchain.Middleware {
	return OnPrefix(prefix, mchain.HandlerFunc(h))
}

func Mount(prefix string, h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			done, err := utils.OnStrippedPrefix(prefix, h, w, r)
			if done {
				return err
			}
			return next.ServeHTTP(w, r)
		}
		return mchain.HandlerFunc(f)
	}
	return hh
}

func MountRedirectToSlashed(prefix string, h mchain.Handler) mchain.Middleware {
	hh := func(next mchain.Handler) mchain.Handler {
		f := func(w http.ResponseWriter, r *http.Request) error {
			done, err := utils.OnStrippedPrefixAndRedirectToSlash(prefix, h, w, r)
			if done {
				return err
			}
			return next.ServeHTTP(w, r)
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
