# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=proxyservice
BINARY_UNIX=$(BINARY_NAME)_unix

test: 
	# $(GOBUILD) -o $(BINARY_NAME) -v
	# $(GOTEST) -v ./...
	docker-compose build
