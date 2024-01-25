package server

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/sharedlogger"
	"fmt"
	"io"
	"net/http"
)

const LOCALHOST string = "http://localhost"
const FORWARDED_FOR string = "X-Forwarded-For"
const FORWARDED string = "Forwarded"

type ProxyHandler struct {
	proxyMap map[string]int
	logger   *sharedlogger.SharedLogger
}

func ProxyHandlerInit(conf *config.ReverseProxy, logger *sharedlogger.SharedLogger) *ProxyHandler {
	pmap := make(map[string]int)
	for _, rprox := range conf.Routes {
		pmap[rprox.From] = rprox.To
		logger.Finfo("Registered proxy from [%s] -> [:%d]", rprox.From, rprox.To)
	}
	return &ProxyHandler{
		proxyMap: pmap,
		logger:   logger,
	}
}

// Should we proxy the request
func (p *ProxyHandler) IsProxied(r *http.Request) bool {
	_, present := p.proxyMap[r.URL.Path]
	return present
}

func (p *ProxyHandler) handleHeaders(r *http.Request, target string) {
	for k, _ := range r.Header {
		switch k {
		case "Connection":
			r.Header[k] = []string{"Upgrade"}
		//We delete this header for security reasons
		case "Referer":
			delete(r.Header, k)
		}
	}
	r.Header.Add(FORWARDED_FOR, r.RemoteAddr)
	r.Header.Add(FORWARDED, r.RemoteAddr)
}

// Handle the response from the proxy call and return it to the caller
func (p *ProxyHandler) handleResponse(r *http.Response, w http.ResponseWriter, errorCallback func(w http.ResponseWriter, uuid string), uuid string) {
	//todo construct response from the resp and we should be good to go
	res, err := io.ReadAll(r.Body)
	if err != nil {
		errorCallback(w, uuid)
		return
	}
	for key, val := range r.Header {
		for _, v := range val {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(r.StatusCode)
	w.Write(res)
}

// Proxy the current request to the given port on the localhost
func (p *ProxyHandler) ProxyRequest(r *http.Request, w http.ResponseWriter, errorCallback func(w http.ResponseWriter, uuid string), uuid string) {
	port, found := p.proxyMap[r.URL.Path]
	if !found {
		p.logger.Ferror("Not found proxy port for url %s.", r.URL.Path)
		return
	}
	//Construction of the redirection url path
	url := fmt.Sprintf("%s:%d/%s", LOCALHOST, port, r.URL.Path)
	p.handleHeaders(r, url)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		errorCallback(w, uuid)
	}
	p.handleResponse(resp, w, errorCallback, uuid)
}
