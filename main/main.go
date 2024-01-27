package main

import (
	"cz/jakvitov/webserv/config"
	"cz/jakvitov/webserv/err"
	"cz/jakvitov/webserv/server"
	"cz/jakvitov/webserv/static"
	"fmt"
	"os"
	"sync"
)

const HELP string = "--help"

func main() {
	if len(os.Args) != 2 {
		err.ErrorPrompt(fmt.Sprintf("Invalid number of arguments, required 1, found %d\n", len(os.Args)-1))
		return
	}

	path := os.Args[1]

	if path == HELP {
		static.HelpMenu()
		return
	}

	cnf, readErr := config.ReadAndVerify(path)
	if readErr != nil {
		err.ErrorPrompt(fmt.Sprintf("Cannot setup config. %s\n", readErr.Error()))
		return
	}
	//Wait group for the server to finish
	terminated := new(sync.WaitGroup)

	srv := server.ServerInit(cnf, terminated)
	wg := srv.StartListening()
	//Wait until the server starts up
	wg.Wait()

	//Block until the server is terminated
	terminated.Wait()
}
