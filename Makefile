GOCMD=go
GOBUILD=$(GOCMD) build -v 
GOHOSTOS=$(strip $(shell $(GOCMD) env get GOHOSTOS))

TAG ?= $(shell git describe --tags)
COMMIT ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date -u +%m/%d/%Y)

# Active module mode, as we use go modules to manage dependencies
export GO111MODULE=on

BUPKIS=bin/bupkis
BUPKIS_DAR=bin/bupkis-darwin

all: format bupkis bupkisdar

clean:
	rm -rf ${BUPKIS} 
bupkisdar:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${BUPKIS_DAR} github.com/zawachte-msft/bupkis
bupkis:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -ldflags "-X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" -o ${BUPKIS} github.com/zawachte-msft/bupkis

.PHONY: vendor
vendor:
	go mod tidy

format:
	gofmt -s -w cmd/ 