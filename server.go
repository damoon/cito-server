package cito

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunServer(addr string) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", fetch)
	mux.HandleFunc("/healthz", healthz)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  1800 * time.Second,
		WriteTimeout: 1800 * time.Second,
	}

	// wait for an exit signal
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-stop
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Fatalf("server shutdown failed: %s\n", err)
		}
	}()

	// serve requests
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
