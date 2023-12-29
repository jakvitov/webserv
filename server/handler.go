package server

import (
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"net/http"
	"os"
)

const BAD_REQUEST string = "Bad request"
const NOT_FOUND string = "Not found"
const OK string = "Ok"

// Holds all the muxes and routes each request to its proper handler
type HttpRequestHandler struct {
	mux    *http.ServeMux
	logger *sharedlogger.SharedLogger
	root   string
}

func (h *HttpRequestHandler) badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(505)
	w.Write([]byte(BAD_REQUEST))
	h.logger.LogHttpResponse(505, r)
}

func (h *HttpRequestHandler) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(NOT_FOUND))
	h.logger.LogHttpResponse(404, r)
}

func (h *HttpRequestHandler) ok(w http.ResponseWriter, r *http.Request, response []byte) {
	w.WriteHeader(200)
	w.Write(response)
	h.logger.LogHttpResponse(200, r)
}

// Request handlers for server requests
func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.LogHttpRequest(r)
	path := r.URL.Path
	if path == "/" {
		path = "index.html"
	}
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", h.root, path))
	if err != nil {
		h.notFound(w, r)
		return
	}
	h.ok(w, r, file)
}

// We map routes to individual handlers in a mux
func (h *HttpRequestHandler) registerHandlers() {
	h.mux.Handle("/", h)
}

func HttpRequestHandlerInit(il *sharedlogger.SharedLogger, root string) *HttpRequestHandler {
	mux := http.NewServeMux()
	res := &HttpRequestHandler{
		mux:    mux,
		logger: il,
		root:   root,
	}
	//Register handlers
	res.registerHandlers()
	return res
}
