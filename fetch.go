package cito

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/minio/minio-go"
)

func fetch(hc *http.Client, mc *minio.Client, bucket string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := ""
		switch {
		case strings.HasPrefix(r.URL.Path, "/http/"):
			u = fmt.Sprintf("http://%s", strings.TrimPrefix(r.URL.Path, "/http/"))
		case strings.HasPrefix(r.URL.Path, "/https/"):
			u = fmt.Sprintf("https://%s", strings.TrimPrefix(r.URL.Path, "/https/"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("unknown protocol %s\n", r.URL.Path)
			return
		}

		o := strings.TrimPrefix(r.URL.Path, "/")

		exists, err := objectExists(r.Context(), mc, bucket, o)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("s3 object %s lookup failed: %s\n", o, err)
			return
		}
		if !exists {
			if os.Getenv("DEBUG") == "1" {
				log.Printf("fetching %s\n", u)
			}

			resp, err := hc.Get(u)
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

			// TODO: use local directory as cache, limit local cache to x MiB

			// TODO: upload to s3 after serving request

			// TODO: gc s3

			_, err = mc.PutObjectWithContext(r.Context(), bucket, o, resp.Body, resp.ContentLength, minio.PutObjectOptions{})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("uploading %s failed: %s\n", o, err)
				return
			}
		}

		// TODO: add content type
		obj, err := mc.GetObjectWithContext(r.Context(), bucket, o, minio.GetObjectOptions{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("s3 streaming %s failed: %s\n", o, err)
			return
		}
		defer obj.Close()

		if _, err := io.Copy(w, obj); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("failed to send result: %s\n", err)
		}
	}
}
