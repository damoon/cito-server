package cito

import (
	"log"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {

	// TODO: fail during shutdown, allow time to shut down

	// TODO: check minio availability

	// TODO: skip minio if last use happened after last healthcheck

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("failed to send 'OK': %s\n", err)
	}
}
