package hostrouter

import (
	"net/http"
)

type HostRouter struct {
	Items     interface{}
	Threshold int
}

type HostRouterItem struct {
	name    string
	handler http.Handler
}

func New() *HostRouter {
	return &HostRouter{Threshold: 7}
}

func (h *HostRouter) Build(notFoundHandler http.Handler) http.Handler {
	if items, ok := h.Items.(map[string]http.Handler); ok {
		h := func(w http.ResponseWriter, r *http.Request) {
			if handler, ok := items[r.URL.Hostname()]; ok {
				handler.ServeHTTP(w, r)
			} else {
				notFoundHandler.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(h)
	}
	items := h.Items.([]HostRouterItem)
	hx := func(w http.ResponseWriter, r *http.Request) {
		for _, x := range items {
			if x.name == r.URL.Hostname() {
				x.handler.ServeHTTP(w, r)
				return
			}
			notFoundHandler.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(hx)
}

func (h *HostRouter) resolveContainer() {
	switch item := h.Items.(type) {
	case map[string]http.Handler:
		if len(item) < h.Threshold+1 {
			s := make([]HostRouterItem, 0, len(item)+1)
			for k, v := range item {
				s = append(s, HostRouterItem{name: k, handler: v})
			}
			h.Items = s
		}
	case []HostRouterItem:
		if len(item)+1 > h.Threshold {
			m := make(map[string]http.Handler, len(item)+1)
			for _, x := range item {
				m[x.name] = x.handler
			}
			h.Items = m
		}
	default:
		h.Items = []HostRouterItem{}
	}
}

func (h *HostRouter) Set(host string, handler http.Handler) {
	h.resolveContainer()
	switch item := h.Items.(type) {
	case map[string]http.Handler:
		if handler == nil {
			delete(item, host)
		} else {
			item[host] = handler
		}
	case []HostRouterItem:
		if handler == nil {
			for i, x := range item {
				if x.name == host {
					// Remove item
					h.Items = append(item[:i], item[i+1:]...)
					break
				}
			}
		} else {
			route := HostRouterItem{host, handler}
			for i, x := range item {
				if x.name == host {
					// Replace item
					item[i] = route
					return
				}
			}
			h.Items = append(item, route)
		}
	}
}
