package cito

import (
	"log"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("failed to send 'OK': %s\n", err)
	}
}
