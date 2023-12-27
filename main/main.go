package main

import (
	"cz/jakvitov/webserv/config"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1]
	cnf, err := config.ReadConfig(args)
	if err != nil {
		fmt.Printf("Error while starting the application: \n\t->%s\n", err.Error())
	}

}
