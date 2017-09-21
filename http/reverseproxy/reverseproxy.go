package reverseproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewURLProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		appendQueryIfNeeded(req, target.RawQuery)
		clearUserAgentIfNotValid(req)
	}
	return &httputil.ReverseProxy{Director: director}
}

func NewHostProxy(host string, forceHttp bool, replaceHostHeader bool) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		if forceHttp {
			req.URL.Scheme = "http"
		}
		if replaceHostHeader {
			req.Header.Set("X-Forwarded-Host", req.Host)
			req.Host = host
		} else {
			req.URL.Host = req.Host
		}
		clearUserAgentIfNotValid(req)
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// WARN: Is this really needed?
func clearUserAgentIfNotValid(req *http.Request) {
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
}

func appendQueryIfNeeded(req *http.Request, targetRawQuery string) {
	if targetRawQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetRawQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetRawQuery + "&" + req.URL.RawQuery
	}
}
