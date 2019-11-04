package cito

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go"
)

func fetch(hc *http.Client, mc *minio.Client, bucket string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		remoteUrl, path := extractRequest(r)
		log.Printf("request for file: %s\n", remoteUrl)

		found, err := objectExists(r.Context(), mc, bucket, path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("looking up file: %v", err)
			return
		}
		if !found {
			err = archivePackage(r.Context(), hc, remoteUrl, mc, bucket, path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("storing package: %v", err)
				return
			}
		}

		// redirect to mino
		signedUrl, err := mc.PresignedGetObject(bucket, path, 5*time.Minute, url.Values{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("signing redirect to minio: %v", err)
		}

		http.Redirect(w, r, signedUrl.String(), http.StatusTemporaryRedirect)
	}
}

func extractRequest(r *http.Request) (string, string) {
	url := ""
	switch {
	case strings.HasPrefix(r.URL.Path, "/http/"):
		url = fmt.Sprintf("http://%s", strings.TrimPrefix(r.URL.Path, "/http/"))
	case strings.HasPrefix(r.URL.Path, "/https/"):
		url = fmt.Sprintf("https://%s", strings.TrimPrefix(r.URL.Path, "/https/"))
	}

	minioPath := strings.TrimPrefix(r.URL.Path, "/")

	return url, minioPath
}

func archivePackage(ctx context.Context, hc *http.Client, remoteUrl string, mc *minio.Client, bucket, path string) error {
	tmp, err := ioutil.TempFile(os.TempDir(), "cito")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	// download
	err = downloadFromRemote(ctx, hc, remoteUrl, tmp)
	if err != nil {
		return fmt.Errorf("downloading package: %v", err)
	}
	_, err = tmp.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seeking temp file: %v", err)
	}

	// upload
	err = uploadToMinio(tmp, mc, bucket, path)
	if err != nil {
		return fmt.Errorf("package upload to minio: %v", err)
	}

	return nil
}

func downloadFromRemote(ctx context.Context, hc *http.Client, url string, w io.Writer) error {

	if os.Getenv("DEBUG") == "1" {
		log.Printf("downloading from remote %s", url)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recieved status code %d", resp.StatusCode)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func uploadToMinio(tmp *os.File, mc *minio.Client, bucket, path string) error {
	if os.Getenv("DEBUG") == "1" {
		log.Printf("uploading to minio %s/%s\n", bucket, path)
	}

	_, err := tmp.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seeking temp file: %v", err)
	}
	fi, err := tmp.Stat()
	if err != nil {
		return fmt.Errorf("stat temp file: %v", err)
	}

	_, err = mc.PutObject(bucket, path, tmp, fi.Size(), minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("uploading %s/%s failed: %s", bucket, path, err)
	}

	return nil
}
