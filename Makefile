REPO = docker.clevalert.com
NAME = clevalert-daemon
ARCH = x86_64

VERSION = $(or $(shell git log -1 --pretty='format:%cd' --date=format:'%y.%m.%d.%H%M'),0.0.0)
TEMPDIR = $(shell mktemp -d)
PROJECT_PATH = $(shell pwd)
GIT_HASH = $(shell git rev-parse --short HEAD)
RPM_PATH = $(shell pwd)/build/rpmbuild/
RPM_FILE = $(RPM_PATH)$(ARCH)/$(NAME)-$(VERSION)-$(GIT_HASH).$(ARCH).rpm
TEST_REPORT_FILE = report.xml

.PHONY: build
build:
	env VERSION=$(VERSION) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-tags static_all \
	-ldflags "-s -X $(NAME)/internal/pkg/configs.Version=${VERSION}" \
	-mod vendor \
	-a -o $(NAME) cmd/run.go

.PHONY: build_docker_image
build_docker_image:
	docker build --target build --tag $(VERSION) \
		--build-arg NAME=$(NAME) \
		--build-arg VERSION=$(VERSION) \
		--file build/Dockerfile .

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: validate
validate:
	go run -race  cmd/run.go -config configs/config.toml -validate

.PHONY: run
run:
	sudo go run -race cmd/run.go -config configs/config.toml

.PHONY: pg_tunnel
pg_tunnel:
	ssh -N -L 5432:localhost:5432 root@77.223.97.227
