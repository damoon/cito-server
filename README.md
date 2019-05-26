# Cito Server

This accept http(s) requests over http to be able to cache (https hard coded) composer packages.

## testing

```
export GO111MODULE=on

docker run --rm -p 9000:9000 --name minio \
  -e "MINIO_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE" \
  -e "MINIO_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" \
  minio/minio server /data

go run ./cmd/main.go -endpoint=127.0.0.1:9000 -useSSL=false -accessKeyID=AKIAIOSFODNN7EXAMPLE -secretAccessKey=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

## TODO

```
Makefile
golangci-lint
docker-compose / gitlab ci test
Dockerfile
kustomize
```
