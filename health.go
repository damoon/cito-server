package cito

import (
	"context"
	"log"
	"net/http"

	"github.com/minio/minio-go"
)

func newHealth(mc *minio.Client, bucket string) http.HandlerFunc {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			// TODO: fail during shutdown, allow time to shut down

			// TODO: skip minio if last use happened after last healthcheck

			if minioUnavilable(r.Context(), mc, bucket) {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}

			w.WriteHeader(http.StatusOK)
		})
}

func minioUnavilable(ctx context.Context, mc *minio.Client, bucket string) bool {
	exists, err := objectExists(ctx, mc, bucket, "/")
	if err != nil {
		log.Printf("failed stat / in bucket %s: %s\n", bucket, err)
		return false
	}
	if !exists {
		return false
	}
	return true
}
