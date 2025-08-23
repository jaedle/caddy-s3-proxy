package internal

import (
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type handler struct {
}

func New() caddyhttp.MiddlewareHandler {
	return &handler{}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, _ *http.Request, _ caddyhttp.Handler) error {
	w.WriteHeader(http.StatusNotFound)
	return nil
}
