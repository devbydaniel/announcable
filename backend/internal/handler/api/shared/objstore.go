package shared

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5"
)

// HandleObjStore proxies requests to object storage
func (h *Handlers) HandleObjStore(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleObjStore")
	// Internal service URL (accessible within Docker network)
	target, err := url.Parse("http://objstorage:9000")
	if err != nil {
		panic(err)
	}

	// Create proxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify the director to ensure proper host header
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host // This is important for internal Docker DNS resolution
	}

	remainingPath := chi.URLParam(r, "*")
	r.URL.Host = target.Host
	r.URL.Scheme = target.Scheme
	r.URL.Path = "/" + remainingPath
	proxy.ServeHTTP(w, r)
}
