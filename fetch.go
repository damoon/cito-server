package cito

import (
	"net/http"
)

func fetch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
