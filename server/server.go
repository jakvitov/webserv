package server

import (
	"context"
	"cz/jakvitov/webserv/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const INFO_LOG_PREFIX string = "WEBSERV_SERVER [INFO]:"
const ERROR_LOG_PREFIX string = "WEBSERV_SERVER [ERROR]:"

// Crentral struct holding info about all the http listeners and config
type Server struct {
	cnf         *config.WebserverConfig
	httpServers []*http.Server
	infoLogger  *log.Logger
	errorLogger *log.Logger
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
		infoLogger:  log.New(io.MultiWriter(os.Stdout, writer), INFO_LOG_PREFIX, log.Ltime),
		errorLogger: log.New(io.MultiWriter(os.Stderr, writer), ERROR_LOG_PREFIX, log.Ltime),
	}
	return srv
}

func (s *Server) StartListening() {
	for _, srv := range s.httpServers {
		go func(s *Server, srv *http.Server) {
			s.infoLogger.Printf("Starting listener on port [%s]\n", srv.Addr)
			log.Fatal(srv.ListenAndServe())
			defer s.errorLogger.Fatal(srv.Shutdown(context.Background()))
		}(s, srv)
	}
}

func (s *Server) Shutdown() {
	for _, srv := range s.httpServers {
		s.infoLogger.Printf("Trying mercifull shutdown on listener,  port [%s]\n", srv.Addr)
		s.errorLogger.Fatal(srv.Shutdown(context.Background()))
	}
}
