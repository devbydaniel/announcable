package shared

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/devbydaniel/announcable/config"
	"github.com/go-chi/chi/v5"
)

// HandleObjStore proxies requests to object storage
func (h *Handlers) HandleObjStore(w http.ResponseWriter, r *http.Request) {
	h.Log.Trace().Msg("HandleObjStore")

	// Build target URL from config
	// Always use HTTP for internal Docker communication (SSL is handled by nginx externally)
	cfg := config.New()
	target, err := url.Parse("http://" + cfg.ObjStorage.Endpoint)
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
