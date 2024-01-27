package echoserver

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

// Basic echo server used for testing purposes
// Returns channel to stop the server
func RunEchoServer(port int, startup *sync.WaitGroup, res chan bool, sh chan bool) {
	startup.Add(1)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
	fmt.Printf("Server is listening on port %d...\n", port)
	startup.Done()
	server := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
	go func() {
		<-res
		server.Close()
		sh <- true
		fmt.Println("Closed test echo server.")
	}()
}
