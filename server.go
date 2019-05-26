package cito

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minio/minio-go"
)

func RunServer(httpClient *http.Client, mc *minio.Client, bucket, addr string) {

	// TODO: use http.TimeoutHandler
	mux := http.NewServeMux()
	mux.HandleFunc("/", fetch(httpClient, mc, bucket))
	mux.HandleFunc("/healthz", healthz)

	// TODO: separate user (8080) and admin endpoint (8081)

	// TODO: add USE, RED and golang metrics

	// TODO: add profiling https://matoski.com/article/golang-profiling-flamegraphs/

	// TODO: add debuging https://github.com/Microsoft/vscode-go/wiki/Debugging-Go-code-using-VS-Code

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
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
