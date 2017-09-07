package app

import "github.com/rsms/gotalk"

func apiHandlers(context *AppContext) *gotalk.Handlers {
	h := gotalk.NewHandlers()

	h.Handle("hello", func() (string, error) {
		return "Hello there!", nil
	})

	h.Handle("echo", func(in string) (string, error) {
		return in, nil
	})

	return h
}
