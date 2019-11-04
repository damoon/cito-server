package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/damoon/cito-server"
	"github.com/minio/minio-go"
	"github.com/pkg/profile"
)

func main() {

	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	serviceAddr := flag.String("service-address", ":8080", "service server address, ':8080'")
	adminAddr := flag.String("admin-address", ":8081", "admin server address, ':8081'")
	shutdownDelay := flag.Duration("shutdown-delay", 30*time.Second, "time in seconds to allow ingress controllers to change routes")
	endpoint := flag.String("endpoint", "", "s3 endpoint")
	accessKeyID := flag.String("accessKeyID", "", "s3 accessKeyID")
	secretAccessKey := flag.String("secretAccessKey", "", "s3 secretAccessKey")
	useSSL := flag.Bool("useSSL", true, "s3 uses https")
	bucket := flag.String("bucket", "cito", "s3 bucket name")
	location := flag.String("location", "us-east-1", "s3 bucket location")

	flag.Parse()

	log.Printf("service server listens on: %s\n", *serviceAddr)
	log.Printf("admin server listens on: %s\n", *adminAddr)

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	var httpClient = &http.Client{
		Timeout:   time.Second * 20,
		Transport: transport,
	}

	minioClient, err := minio.New(*endpoint, *accessKeyID, *secretAccessKey, *useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	minioClient.SetCustomTransport(transport)

	go func() {
		for {
			err = cito.EnsureBucket(minioClient, *bucket, *location)
			if err != nil {
				log.Printf("failed to ensure bucket exists: %v\n", err)
				time.Sleep(5 * time.Second)
				continue
			}
			return
		}
	}()

	cito.RunServer(httpClient, minioClient, *bucket, *adminAddr, *serviceAddr, *shutdownDelay)
}
