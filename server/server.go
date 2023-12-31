package server

import (
	"context"
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/sharedlogger"
	"cz/jakvitov/webserv/static"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const INFO_LOG_PREFIX string = "WEBSERV_SERVER [INFO]:"
const ERROR_LOG_PREFIX string = "WEBSERV_SERVER [ERROR]:"

// Crentral struct holding info about all the http listeners and config
type Server struct {
	cnf         *config.WebserverConfig
	httpServers []*http.Server
	logger      *sharedlogger.SharedLogger
}

func initHttpServers(cnf *config.WebserverConfig, logger *sharedlogger.SharedLogger) []*http.Server {
	res := make([]*http.Server, len(cnf.Ports))
	for i, port := range cnf.Ports {
		res[i] = &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        HttpRequestHandlerInit(logger, cnf.RootDir),
			ReadTimeout:    1 * time.Second,
			WriteTimeout:   1 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	}
	return res
}

func ServerInit(inputCnf *config.WebserverConfig) *Server {
	//Open as create or append
	writer, err := os.OpenFile(inputCnf.LogDest, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("Error opening log file [%s], creating one instead.\n", inputCnf.LogDest)
		writerC, err := os.Create(inputCnf.LogDest)
		if err != nil {
			panic(fmt.Sprintf("Error while creating a log file: [%s]\n", err.Error()))
		}
		writer = writerC
	}

	lg := sharedlogger.SharedLoggerInit(writer)
	srv := &Server{
		cnf:         inputCnf,
		logger:      lg,
		httpServers: initHttpServers(inputCnf, lg),
	}
	return srv
}

func (s *Server) StartListening(wg *sync.WaitGroup) {
	s.ListenForSigterm()
	static.PrintBannerDecoration(s.logger)
	for _, srv := range s.httpServers {
		wg.Add(1)
		go func(s *Server, srv *http.Server) {
			s.logger.Finfo("Starting listener on port [%s]\n", srv.Addr)
			err := srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				s.logger.Fatal(err.Error())
			}
			defer wg.Done()
		}(s, srv)
	}
}

// Force quit all listening servers
func (s *Server) Shutdown() {
	for _, srv := range s.httpServers {
		if err := srv.Shutdown(context.Background()); err != nil {
			s.logger.Error(err.Error())
		}
	}
}

// Listends for sigterm system signal and tries to gracefully shutdown afterwards
func (s *Server) ListenForSigterm() {
	//Channel listening to sigterm signal
	sigNotif := make(chan os.Signal, 1)
	signal.Notify(sigNotif,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		//Listen for termination signal
		sig := <-sigNotif
		s.logger.Finfo("Recieved %s signal. Attempting gracefull shutdown.", sig.String())
		s.Shutdown()
	}()
}
