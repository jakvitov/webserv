package server

import (
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var URL_BLACKLIST []string = []string{"../"}

const BAD_REQUEST string = "Bad request"
const NOT_FOUND string = "Not found"

// Holds all the muxes and routes each request to its proper handler
type HttpRequestHandler struct {
	mux    *http.ServeMux
	logger *sharedlogger.SharedLogger
	root   string
}

func (h *HttpRequestHandler) badRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(505)
	w.Write([]byte(BAD_REQUEST))
	h.logger.Info("Returning 505")
	return
}

func (h *HttpRequestHandler) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte(NOT_FOUND))
	h.logger.Info("Returning 404")
	return
}

// Request handlers for server requests
func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Check for items on the black list and return bad request if present
	for _, item := range URL_BLACKLIST {
		if strings.Contains(r.URL.RequestURI(), item) {
			h.badRequest(w, r)
			return
		}
	}
	file, err := os.ReadFile(fmt.Sprintf("%s%s", h.root, "index.html"))
	if err != nil {
		h.notFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

// We map routes to individual handlers in a mux
func (h *HttpRequestHandler) registerHandlers() {
	h.mux.Handle("/", h)
}

func HttpRequestHandlerInit(il *sharedlogger.SharedLogger) *HttpRequestHandler {
	mux := http.NewServeMux()
	res := &HttpRequestHandler{
		mux:    mux,
		logger: il,
	}
	//Register handlers
	res.registerHandlers()
	return res
}
