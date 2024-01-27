package server

import (
	"context"
	"crypto/tls"
	"cz/jakvitov/webserv/cache"
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

// Crentral struct holding info about all the http listeners and config
type Server struct {
	cnf         *config.Config
	httpServer  *http.Server
	httpsServer *http.Server
	//Lookup map for reverse proxy fast lookups
	logger *sharedlogger.SharedLogger
	termWg *sync.WaitGroup
}

func initHttpServer(cnf *config.Config, logger *sharedlogger.SharedLogger) *http.Server {
	logger.Finfo("Creating http server for port [%d]", cnf.Ports.HttpPort)
	cache := cache.CacheInit(cnf.Handler.MaxCacheSize, cnf.Handler.CacheEnabled, logger)
	if cnf.Handler.CacheEnabled {
		logger.Finfo("Created server cache with max size of [%d] bytes.", cnf.Handler.MaxCacheSize)
	}

	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cnf.Ports.HttpPort),
		Handler:        HttpRequestHandlerInit(logger, cnf, cache),
		ReadTimeout:    time.Duration(cnf.Handler.ReadTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(cnf.Handler.WriteTimeout) * time.Millisecond,
		MaxHeaderBytes: cnf.Handler.MaxHeaderBytes,
	}
}

// If the config contains https port, start a https server, if not, don't
// Cache and handler should be already ready from the http server
func initHttpsServer(cnf *config.Config, logger *sharedlogger.SharedLogger, handler *http.Handler) *http.Server {
	if cnf.Ports.HttpsPort == 0 {
		return nil
	}
	if cnf.Handler.CacheEnabled {
		logger.Finfo("Using server cache of [%d] bytes for https.", cnf.Handler.MaxCacheSize)
	}

	certif, err := static.LoadTlsCerfificate(cnf.Security.CertPath, cnf.Security.PrivateKeyPath, logger)
	if err != nil {
		logger.Error("Certificate cannot be loaded. HTTPS server won't be launched!")
		return nil
	}
	return &http.Server{
		Addr:           fmt.Sprintf(":%d", cnf.Ports.HttpsPort),
		Handler:        *handler,
		ReadTimeout:    time.Duration(cnf.Handler.ReadTimeout) * time.Millisecond,
		WriteTimeout:   time.Duration(cnf.Handler.WriteTimeout) * time.Millisecond,
		MaxHeaderBytes: cnf.Handler.MaxHeaderBytes,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*certif},
		},
	}
}

// Init the server with the termination wait group and input config
func ServerInit(inputCnf *config.Config, terminateWg *sync.WaitGroup) *Server {
	var lg *sharedlogger.SharedLogger
	if inputCnf.Logger.OutputToFile {
		//Open as create or append
		writer, err := os.OpenFile(inputCnf.Logger.OutputFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Printf("Error opening log file [%s], creating one instead.", inputCnf.Logger.OutputFile)
			writerC, err := os.Create(inputCnf.Logger.OutputFile)
			if err != nil {
				panic(fmt.Sprintf("Error while creating a log file: [%s]", err.Error()))
			}
			writer = writerC
		}
		lg = sharedlogger.SharedLoggerInit(writer, inputCnf)
	} else {
		//User chose not to log into a file => only std
		//WE pass nil as file stream ptr and indicate to logger, that we want only std
		lg = sharedlogger.SharedLoggerInit(nil, inputCnf)
	}
	srv := &Server{
		cnf:        inputCnf,
		logger:     lg,
		httpServer: initHttpServer(inputCnf, lg),
		//Waitgroup to be Done() when the server is shut down
		termWg: terminateWg,
	}
	srv.httpsServer = initHttpsServer(inputCnf, lg, &srv.httpServer.Handler)
	return srv
}

// Returns wait group, that is Done as soon as the server is listening on the give port and finished setup
func (s *Server) StartListening() *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	s.ListenForSigterm()
	static.PrintBannerDecoration(s.logger)
	wg.Add(1)
	s.termWg.Add(1)
	//Start http server
	go func() {
		s.logger.Finfo("Starting http listener on port [:%d]", s.cnf.Ports.HttpPort)
		wg.Done()
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal(err.Error())
		}
	}()

	//No need to init https server - it is disabled
	if s.httpsServer == nil {
		return wg
	}

	//Enabled https -> we need to initialise it
	wg.Add(1)
	s.termWg.Add(1)
	//Start https server
	go func() {
		s.logger.Finfo("Starting https listener on port [:%d]", s.cnf.Ports.HttpsPort)
		wg.Done()
		err := s.httpsServer.ListenAndServeTLS(s.cnf.Security.CertPath, s.cnf.Security.PrivateKeyPath)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Ferror("Error while opening the TLS server [%s]\n", err.Error())
		}
	}()
	return wg
}

// Force quit all listening servers
func (s *Server) Shutdown() {
	//Free the wait groups after shutdown
	defer s.termWg.Done()
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		s.logger.Error(err.Error())
	}
	//We do not have https server
	if s.httpsServer == nil {
		return
	}
	//We have to shutdown the https as well
	defer s.termWg.Done()
	if err := s.httpsServer.Shutdown(context.Background()); err != nil {
		s.logger.Error(err.Error())
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
