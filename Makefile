
include ./etc/help.mk

UID := $(shell id -u)
GID := $(shell id -g)

PHP = docker run --rm -ti \
	-e COMPOSER_HOME=/composer-home \
	-v $(PWD)/test/composer-home:/composer-home \
	-v $(PWD)/test/app:/app \
	-e CITO_SERVER=http://127.0.0.1:8080 \
	-e CITO_DEBUG= \
	--network=host \
	--entrypoint=bash \
	--user $(UID):$(GID) \
	cito-test-docker-image

test/docker/image: test/docker/Dockerfile
	docker build -t cito-test-docker-image ./test/docker
	touch test/docker/image

lint: ##@qa run linting for golang.
	golangci-lint run --enable-all ./...

.PHONY: minio
minio: ##@development Start minio server (port:9000, user:minio, secret:minio123).
	docker run --rm \
	-p 9000:9000 \
	-e MINIO_ACCESS_KEY=minio \
	-e MINIO_SECRET_KEY=minio123 \
	minio/minio \
	server /data
#	-e MINIO_HTTP_TRACE=/dev/stdout \

start: ##@development Start cito server (port:8080, admin port:8081).
	DEBUG=0 air -c air.conf

.PHONY: php
php: test/docker/image ##@development Open a command line interface with PHP & composer installed.
	$(PHP)

.PHONY: test
test: test/docker/image ##@qa Check installing takes less then 1s.
	rm -rf test/app/vendor test/composer-home/cache || true
	$(PHP) -c "time composer install --ignore-platform-reqs --no-scripts --no-autoloader"
#	rm -rf test/app/vendor test/composer-home/cache
#	$(PHP) -c "timeout 1 time composer install --ignore-platform-reqs"
