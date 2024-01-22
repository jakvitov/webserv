package server

import (
	"cz/jakvitov/webserv/cache"
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const BAD_REQUEST string = "Bad request"
const NOT_FOUND string = "Not found"
const OK string = "Ok"

// Holds all the muxes and routes each request to its proper handler
type HttpRequestHandler struct {
	mux    *http.ServeMux
	logger *sharedlogger.SharedLogger
	root   string
	cache  *cache.Cache
}

func (h *HttpRequestHandler) badRequest(w http.ResponseWriter, uuid string) {
	w.WriteHeader(505)
	w.Write([]byte(BAD_REQUEST))
	h.logger.LogHttpResponse(505, uuid)
}

func (h *HttpRequestHandler) notFound(w http.ResponseWriter, uuid string) {
	w.WriteHeader(404)
	w.Write([]byte(NOT_FOUND))
	h.logger.LogHttpResponse(404, uuid)
}

func (h *HttpRequestHandler) ok(w http.ResponseWriter, response []byte, uuid string) {
	w.WriteHeader(200)
	w.Write(response)
	h.logger.LogHttpResponse(200, uuid)
}

// Request handlers for server requests
func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.NewString()
	h.logger.LogHttpRequest(r, uuid)
	path := r.URL.Path
	if path == "/" {
		path = "index.html"
	}
	file, err := h.cache.Get(fmt.Sprintf("%s/%s", h.root, path))
	if err != nil {
		h.notFound(w, uuid)
		return
	}
	h.ok(w, file, uuid)
}

// We map routes to individual handlers in a mux
func (h *HttpRequestHandler) registerHandlers() {
	h.mux.Handle("/", h)
}

func HttpRequestHandlerInit(il *sharedlogger.SharedLogger, root string, cache *cache.Cache) *HttpRequestHandler {
	mux := http.NewServeMux()
	res := &HttpRequestHandler{
		mux:    mux,
		logger: il,
		root:   root,
		cache:  cache,
	}
	//Register handlers
	res.registerHandlers()
	return res
}
