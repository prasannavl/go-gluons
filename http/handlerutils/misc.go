package handlerutils

import (
	"io"
	"net/http"
	"os"

	"github.com/prasannavl/mchain"
)

func Nop(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NopHandler() mchain.Handler {
	return mchain.HandlerFunc(Nop)
}

func SendFileHandler(filePath string, status int) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()
		w.WriteHeader(status)
		_, err = io.Copy(w, f)
		return err
	}
	return mchain.HandlerFunc(f)
}

func SendStatusHandler(status int) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(status)
		return nil
	}
	return mchain.HandlerFunc(f)
}

func InvokeFunctionHandler(fn func()) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		fn()
		return nil
	}
	return mchain.HandlerFunc(f)
}

func FunctionResultHandler(fn func() error) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		return fn()
	}
	return mchain.HandlerFunc(f)
}

func InvokeFunctionStatusHandler(fn func(), status int) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(status)
		fn()
		return nil
	}
	return mchain.HandlerFunc(f)
}

func FunctionResultStatusHandler(fn func() error, status int) mchain.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(status)
		return fn()
	}
	return mchain.HandlerFunc(f)
}
