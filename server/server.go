package server

import (
	"cz/jakvitov/webserv/config"
	"fmt"
	"net/http"
	"time"
)

// Crentral struct holding info about all the http listeners and config
type Server struct {
	cnf         *config.WebserverConfig
	httpServers []*http.Server
}

func initHttpServers(cnf *config.WebserverConfig) []*http.Server {
	res := make([]*http.Server, len(cnf.Ports))
	for i, port := range cnf.Ports {
		res[i] = &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        HttpRequestHandlerInit(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}
	return res
}

func ServerInit(inputCnf *config.WebserverConfig) *Server {
	srv := &Server{
		cnf:         inputCnf,
		httpServers: initHttpServers(inputCnf),
	}
	return srv
}

func (s *Server) StartListening() {
	for _, srv := range s.httpServers {
		srv.ListenAndServe()
	}
}
