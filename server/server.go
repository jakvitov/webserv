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
	cnf         *config.Config
	httpServer  *http.Server
	httpsServer *http.Server
	logger      *sharedlogger.SharedLogger
}

func initHttpServer(cnf *config.Config, logger *sharedlogger.SharedLogger) *http.Server {
	logger.Finfo("Creating http server for port [%d]\n", cnf.Ports.HttpPort)

	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cnf.Ports.HttpPort),
		Handler:        HttpRequestHandlerInit(logger, cnf.Handler.ContentRoot),
		ReadTimeout:    time.Duration(cnf.Handler.ReadTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(cnf.Handler.WriteTimeout) * time.Millisecond,
		MaxHeaderBytes: cnf.Handler.MaxHeaderBytes,
	}
}

func ServerInit(inputCnf *config.Config) *Server {
	var lg *sharedlogger.SharedLogger
	if inputCnf.Logger.OutputToFile {
		//Open as create or append
		writer, err := os.OpenFile(inputCnf.Logger.OutputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Printf("Error opening log file [%s], creating one instead.\n", inputCnf.Logger.OutputFile)
			writerC, err := os.Create(inputCnf.Logger.OutputFile)
			if err != nil {
				panic(fmt.Sprintf("Error while creating a log file: [%s]\n", err.Error()))
			}
			writer = writerC
			lg = sharedlogger.SharedLoggerInit(writer, inputCnf)
		} else {
			//User chose not to log into a file => only std
			//WE pass nil as file stream ptr and indicate to logger, that we want only std
			lg = sharedlogger.SharedLoggerInit(nil, inputCnf)
		}
	}
	srv := &Server{
		cnf:        inputCnf,
		logger:     lg,
		httpServer: initHttpServer(inputCnf, lg),
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
