package main

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/err"
	"cz/jakvitov/webserv/server"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		err.ErrorPrompt(fmt.Sprintf("Invalid number of arguments, required 1, found %d\n", len(os.Args)-1))
		return
	}
	path := os.Args[1]
	cnf, readErr := config.ReadAndVerify(path)
	if readErr != nil {
		err.ErrorPrompt(fmt.Sprintf("Cannot setup config. %s\n", readErr.Error()))
		return
	}
	terminated := make(chan bool, 1)
	srv := server.ServerInit(cnf)
	wg := srv.StartListening(terminated)
	//Wait until the server starts up
	wg.Wait()

	//Block until the server is terminated
	<-terminated
}
