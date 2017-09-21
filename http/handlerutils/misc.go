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
		return SendFromReader(w, f, status)
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

func SendFromReader(w http.ResponseWriter, reader io.Reader, status int) error {
	w.WriteHeader(status)
	_, err := io.Copy(w, reader)
	return err
}

func SendContentFromReader(w http.ResponseWriter, reader io.Reader, contentType string, status int) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", contentType)
	_, err := io.Copy(w, reader)
	return err
}

func SendFromReaderHandler(reader io.Reader, status int) mchain.Handler {
	if reader == nil {
		return SendStatusHandler(status)
	}
	f := func(w http.ResponseWriter, r *http.Request) error {
		return SendFromReader(w, reader, status)
	}
	return mchain.HandlerFunc(f)
}

func SendContentFromReaderHandler(reader io.Reader, contentType string, status int) mchain.Handler {
	if reader == nil {
		return SendStatusHandler(status)
	}
	if contentType == "" {
		return SendFromReaderHandler(reader, status)
	}
	f := func(w http.ResponseWriter, r *http.Request) error {
		return SendContentFromReader(w, reader, contentType, status)
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
