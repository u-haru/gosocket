package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/u-haru/websocks"
)

func main() {
	client := &websocks.Client{}

	flag.StringVar(&client.Host, "h", "0.0.0.0:8000", "Listening Address:Port")
	flag.StringVar(&client.Target, "t", "ws://127.0.0.1:8080/ws", "Target websocket path")
	flag.Parse()

	go func() {
		log.Println("Start serving on " + client.Host + " to " + client.Target)
		if err := client.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	log.Printf("Signal %s received, shutting down...\n", (<-quit).String())
}
