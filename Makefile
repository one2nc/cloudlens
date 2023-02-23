GO_FLAGS   ?=
NAME       := cloudlens
OUTPUT_BIN ?= execs/${NAME}
PACKAGE    := github.com/one2nc/$(NAME)
VERSION    ?= v0.26.10

ßsetup:
	docker-compose up -d

setup-down:
	docker ps -a --format "{{.ID}} {{.Names}}" | grep cloudlens| awk '{print $$1}'| xargs docker stop | xargs docker rm -v

build:
	go build \
	-ldflags "-w -s -X ${PACKAGE}/cmd.version=${VERSION}" \
	-a -tags netgo -o ${OUTPUT_BIN} main.go
	
run: build
	./execs/cloudlens

populate: build
	./execs/cloudlens lspop
