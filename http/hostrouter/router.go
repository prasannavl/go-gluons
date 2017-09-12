package hostrouter

import (
	"net/http"
	"strings"

	"github.com/gobwas/glob"
)

type HostRouter struct {
	Items        interface{}
	Threshold    int
	PatternItems []RouterGlobItem
}

type RouterItem struct {
	host    string
	handler http.Handler
}

type RouterGlobItem struct {
	pattern string
	matcher glob.Glob
	handler http.Handler
}

func New() *HostRouter {
	return &HostRouter{Threshold: 7}
}

func hostName(r *http.Request) string {
	hostname := stripPort(r.Host)
	return strings.ToLower(hostname)
}

func stripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func (h *HostRouter) Build(notFoundHandler http.Handler) http.Handler {
	if items, ok := h.Items.(map[string]http.Handler); ok {
		hh := func(w http.ResponseWriter, r *http.Request) {
			hostname := hostName(r)
			if handler, ok := items[hostname]; ok {
				handler.ServeHTTP(w, r)
				return
			}
			for _, x := range h.PatternItems {
				if x.matcher.Match(hostname) {
					x.handler.ServeHTTP(w, r)
					return
				}
			}
			notFoundHandler.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hh)
	}
	items := h.Items.([]RouterItem)
	hx := func(w http.ResponseWriter, r *http.Request) {
		hostname := hostName(r)
		for _, x := range items {
			if x.host == hostname {
				x.handler.ServeHTTP(w, r)
				return
			}
		}
		for _, x := range h.PatternItems {
			if x.matcher.Match(hostname) {
				x.handler.ServeHTTP(w, r)
				return
			}
		}
		notFoundHandler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hx)
}

func (h *HostRouter) resolveContainer() {
	switch item := h.Items.(type) {
	case map[string]http.Handler:
		if len(item) < h.Threshold+1 {
			s := make([]RouterItem, 0, len(item)+1)
			for k, v := range item {
				s = append(s, RouterItem{host: k, handler: v})
			}
			h.Items = s
		}
	case []RouterItem:
		if len(item)+1 > h.Threshold {
			m := make(map[string]http.Handler, len(item)+1)
			for _, x := range item {
				m[x.host] = x.handler
			}
			h.Items = m
		}
	default:
		h.Items = []RouterItem{}
	}
}

func (h *HostRouter) HandleHost(host string, handler http.Handler) {
	h.resolveContainer()
	switch item := h.Items.(type) {
	case map[string]http.Handler:
		if handler == nil {
			delete(item, host)
		} else {
			item[host] = handler
		}
	case []RouterItem:
		if handler == nil {
			for i, x := range item {
				if x.host == host {
					// Remove item
					h.Items = append(item[:i], item[i+1:]...)
					break
				}
			}
		} else {
			route := RouterItem{host, handler}
			for i, x := range item {
				if x.host == host {
					// Replace item
					item[i] = route
					return
				}
			}
			h.Items = append(item, route)
		}
	}
}

func (h *HostRouter) HandlePattern(globPattern string, handler http.Handler) {
	hasStar := false
	for _, x := range globPattern {
		if x == '*' {
			hasStar = true
			break
		}
	}

	if !hasStar {
		h.HandleHost(globPattern, handler)
		return
	}

	items := h.PatternItems
	// note: ok to copy into x during range
	for i, x := range items {
		if x.pattern == globPattern {
			if handler == nil {
				h.PatternItems = append(items[:i], items[i+1:]...)
				return
			}
			items[i].handler = handler
			return
		}
	}
	h.PatternItems = append(items, RouterGlobItem{
		pattern: globPattern, matcher: glob.MustCompile(globPattern), handler: handler})
}

func (h *HostRouter) Clone() *HostRouter {
	return &HostRouter{
		Items:        h.cloneItems(),
		Threshold:    h.Threshold,
		PatternItems: h.clonePatternItems(),
	}
}

func (h *HostRouter) cloneItems() interface{} {
	if items, ok := h.Items.(map[string]http.Handler); ok {
		newMap := make(map[string]http.Handler, len(items))
		if len(items) < 1 {
			return nil
		}
		for k, v := range items {
			newMap[k] = v
		}
		return newMap
	} else if items, ok := h.Items.([]RouterItem); ok {
		if len(items) < 1 {
			return nil
		}
		newSlice := make([]RouterItem, len(items))
		for i, v := range items {
			newSlice[i] = v
		}
		return newSlice
	}
	return nil
}

func (h *HostRouter) clonePatternItems() []RouterGlobItem {
	var items []RouterGlobItem
	if len(h.PatternItems) < 1 {
		return items
	}
	items = make([]RouterGlobItem, len(h.PatternItems))
	for i, v := range h.PatternItems {
		items[i] = v
	}
	return items
}
