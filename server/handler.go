package server

import "net/http"

// Holds all the muxes and routes each request to its proper handler
type HttpRequestHandler struct {
	mux *http.ServeMux
}

// Request handlers for server requests
func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Under construction"))
}

// We map routes to individual handlers in a mux
func (h *HttpRequestHandler) registerHandlers() {
	h.mux.Handle("/", h)
}

func HttpRequestHandlerInit() *HttpRequestHandler {
	mux := http.NewServeMux()
	res := &HttpRequestHandler{
		mux: mux,
	}
	//Register handlers
	res.registerHandlers()
	return res
}
