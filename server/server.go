package server

import (
	"cz/jakvitov/webserv/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const LOG_PREFIX string = "WEBSERV_SERVER:"

// Crentral struct holding info about all the http listeners and config
type Server struct {
	cnf         *config.WebserverConfig
	httpServers []*http.Server
	logger      *log.Logger
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
	writer, err := os.Open(inputCnf.LogDest)
	if err != nil {
		fmt.Printf("Error opening log file [%s], creating one instead.\n", inputCnf.LogDest)
		writerC, err := os.Create(inputCnf.LogDest)
		if err != nil {
			panic(fmt.Sprintf("Error while creating a log file: [%s]\n", err.Error()))
		}
		writer = writerC
	}

	srv := &Server{
		cnf:         inputCnf,
		httpServers: initHttpServers(inputCnf),
		logger:      log.New(writer, LOG_PREFIX, log.Ltime),
	}
	srv.logger.SetOutput(io.MultiWriter(writer, os.Stdout))
	return srv
}

func (s *Server) StartListening() {
	for _, srv := range s.httpServers {
		s.logger.Printf("Starting listener on port [%s]\n", srv.Addr)
		s.logger.Fatal(srv.ListenAndServe())
	}
}
