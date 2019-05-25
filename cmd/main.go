package main

import (
	"flag"
	"log"

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

	minioClient, err := minio.New(*endpoint, *accessKeyID, *secretAccessKey, *useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	err = cito.EnsureBucket(minioClient, *bucket, *location)
	if err != nil {
		log.Fatalln(err)
	}
	cito.RunServer(minioClient, *bucket, *addr)
}
