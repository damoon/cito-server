package cito

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var netClient = &http.Client{
	Timeout: time.Second * 20,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	},
}

func fetch(w http.ResponseWriter, r *http.Request) {
	u := ""
	switch true {
	case strings.HasPrefix(r.URL.Path, "/http/"):
		u = fmt.Sprintf("http://%s", strings.TrimPrefix(r.URL.Path, "/http/"))
	case strings.HasPrefix(r.URL.Path, "/https/"):
		u = fmt.Sprintf("https://%s", strings.TrimPrefix(r.URL.Path, "/https/"))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("unknown protocol %s\n", r.URL.Path)
		return
	}

	if os.Getenv("DEBUG") == "1" {
		log.Printf("fetching %s\n", u)
	}

	resp, err := netClient.Get(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("fetching %s failed: %s\n", u, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("fetching %s returned status code %d\n", u, resp.StatusCode)
		return
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to send result: %s\n", err)
	}
}
