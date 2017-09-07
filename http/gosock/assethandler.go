package gosock

import (
	"io"
	"net/http"
	"strconv"

	"github.com/prasannavl/mchain"
	"github.com/rsms/gotalk/js"
)

func CreateAssetHandler(assetPath string, apiPath string, enableSourceMap bool) mchain.Handler {
	sourceMapPath := assetPath + ".map"
	f := func(w http.ResponseWriter, r *http.Request) error {
		if r.URL.Path == assetPath {
			serveResource(w, r, func() {
				w.Header()["Content-Type"] = []string{"text/javascript"}
				serveURL := "this.gotalkResponderAt={ws:'" + apiPath + "'};"
				sizeStr := strconv.FormatInt(int64(len(serveURL)+len(gotalkjs.BrowserLibString)), 10)
				w.Header()["Content-Length"] = []string{sizeStr}
				w.WriteHeader(http.StatusOK)
				// Note: w conforms to interface { WriteString(string)(int,error) }
				io.WriteString(w, serveURL)
				io.WriteString(w, gotalkjs.BrowserLibString)
			})
		} else if enableSourceMap && r.URL.Path == sourceMapPath {
			serveResource(w, r, func() {
				w.Header()["Content-Type"] = []string{"application/json"}
				w.Header()["Content-Length"] = []string{strconv.FormatInt(
					int64(len(gotalkjs.BrowserLibSourceMapString)),
					10,
				)}
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, gotalkjs.BrowserLibSourceMapString)
			})
		}
		return nil
	}
	return mchain.HandlerFunc(f)
}

func serveResource(w http.ResponseWriter, r *http.Request, f func()) {
	// serve javascript library
	w.Header()["Cache-Control"] = []string{"public, max-age=300"}
	etag := "\"" + gotalkjs.BrowserLibSHA1Base64 + r.URL.Path + "\""
	w.Header()["ETag"] = []string{etag}
	reqETag := r.Header["If-None-Match"]

	if len(reqETag) != 0 && reqETag[0] == etag {
		w.WriteHeader(http.StatusNotModified)
	} else {
		f()
	}
}
