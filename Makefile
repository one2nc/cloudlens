GO_FLAGS   ?=
NAME       := cloudlens
OUTPUT_BIN ?= execs/${NAME}
PACKAGE    := github.com/one2nc/$(NAME)
VERSION    = v0.1.4
GIT_REV    ?= $(shell git rev-parse --short HEAD)
SOURCE_DATE_EPOCH ?= $(shell date +%s)
ifeq ($(shell uname), Darwin)
DATE       ?= $(shell TZ=UTC date -j -f "%s" ${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
else
DATE       ?= $(shell date -u -d @${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
endif

build:
	go build ${GO_FLAGS} \
	-ldflags "-w -s -X ${PACKAGE}/cmd.version=${VERSION} -X ${PACKAGE}/cmd.commit=${GIT_REV} -X ${PACKAGE}/cmd.date=${DATE}" \
	-o ${OUTPUT_BIN} main.go
	
run: build
	./execs/cloudlens

setup:
	docker-compose up -d 

setup-down:
	docker ps -a --format "{{.ID}} {{.Names}}" | grep cloudlens| awk '{print $$1}'| xargs docker stop | xargs docker rm -v