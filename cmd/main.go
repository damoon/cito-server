package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/damoon/cito-server"
	"github.com/minio/minio-go"
)

func main() {

	addr := flag.String("address", ":8080", "default server address, ':8080'")
	endpoint := flag.String("endpoint", "", "s3 endpoint")
	accessKeyID := flag.String("accessKeyID", "", "s3 accessKeyID")
	secretAccessKey := flag.String("secretAccessKey", "", "s3 secretAccessKey")
	useSSL := flag.Bool("useSSL", true, "s3 uses https")
	bucket := flag.String("bucket", "cito", "s3 bucket name")
	location := flag.String("location", "us-east-1", "s3 bucket location")

	flag.Parse()

	log.Printf("server listens on: %s\n", *addr)

	// TODO: fail if config is missing

	// TODO: use logrus and json logs

	// TODO: add tracing

	var httpClient = &http.Client{
		Timeout: time.Second * 20,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	minioClient, err := minio.New(*endpoint, *accessKeyID, *secretAccessKey, *useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	err = cito.EnsureBucket(minioClient, *bucket, *location)
	if err != nil {
		log.Fatalln(err)
	}
	cito.RunServer(httpClient, minioClient, *bucket, *addr)
}
