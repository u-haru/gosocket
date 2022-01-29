package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/u-haru/websocks"
)

func main() {
	server := &websocks.Server{}

	flag.StringVar(&server.URI, "h", "0.0.0.0:8080/ws", "Listening Address:Port/Path")
	flag.Parse()

	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Connect to "+server.Path()+" with client")
	})
	go func() {
		log.Println("Start serving on " + server.URI)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	log.Printf("Signal %s received, shutting down...\n", (<-quit).String())
}
