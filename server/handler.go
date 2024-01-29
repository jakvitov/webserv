package server

import (
	"cz/jakvitov/webserv/cache"
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const BAD_REQUEST string = "Bad request"
const NOT_FOUND string = "Not found"
const INTERNAL_SERVER_ERROR string = "Internal server error"
const OK string = "Ok"

// Holds all the muxes and routes each request to its proper handler
type HttpRequestHandler struct {
	mux          *http.ServeMux
	logger       *sharedlogger.SharedLogger
	root         string
	cache        *cache.Cache
	proxyHandler *ProxyHandler
}

//lint:ignore U1000
func (h *HttpRequestHandler) internalServerError(w http.ResponseWriter, uuid string) {
	w.WriteHeader(500)
	w.Write([]byte(INTERNAL_SERVER_ERROR))
	h.logger.LogHttpResponse(500, uuid)
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
	//If the request is proxied, handle it
	if h.proxyHandler.IsProxied(r) {
		h.proxyHandler.ProxyRequest(r, w, h.badRequest, uuid)
		return
	}
	//Route to index on base path
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

func HttpRequestHandlerInit(lg *sharedlogger.SharedLogger, cnf *config.Config, cache *cache.Cache) *HttpRequestHandler {
	mux := http.NewServeMux()
	res := &HttpRequestHandler{
		mux:          mux,
		logger:       lg,
		root:         cnf.Handler.ContentRoot,
		cache:        cache,
		proxyHandler: ProxyHandlerInit(&cnf.ReverseProxy, lg),
	}
	//Register handlers
	res.registerHandlers()
	return res
}
