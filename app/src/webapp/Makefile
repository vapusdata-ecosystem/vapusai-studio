# Makefile

#input params
OCI_IMAGE_PATH = asia-south1-docker.pkg.dev/vapusdata-beta/vapusdata-oss/webapp
TAG = local_$(shell whoami)

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOLINT = golangci-lint
GOSEC = gosec
BUFCMD = buf

# Build target
BINARY_NAME = webapp

all: test build

protogen:
	cd ../../../apis && make gosync

build-binary: protogen gosync build-tailwind
	cd cmd && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix -ldflags="-w -s" -o  ../$(BINARY_NAME) 

build-tailwind:
	npx --yes tailwindcss -i static/css/input.css -o static/css/vapus.css --minify

build: build-binary
	docker build --no-cache -t $(OCI_IMAGE_PATH):$(TAG) .
	rm -f $(BINARY_NAME)

release: build
	docker push $(OCI_IMAGE_PATH):$(TAG) 

# build: protogen gosync
# 	goreleaser build --snapshot --clean --single-target --env CONTAINER_REGISTRY=$(CONTAINER_REGISTRY) --env OCI_IMAGE=$(OCI_IMAGE) --env TAG=$(TAG)

# release: build
# 	goreleaser release --clean --snapshot --env CONTAINER_REGISTRY=$(CONTAINER_REGISTRY) --env OCI_IMAGE=$(OCI_IMAGE) --env TAG=$(TAG)



test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

runmain: build-tailwind
	$(GOCMD) run cmd/main.go --conf=$(CONF)

lint:
	$(GOLINT) run

sec:
	$(GOSEC) ./...

.PHONY: all build test clean run

prehook: lint sec test

gosync: 
	cd ../../../internals && make gosync
	$(GOCMD) mod tidy 
	$(GOCMD) mod vendor
goget:
	echo argument is $(lib)
	$(GOGET) $(lib)
	make gosync

goinstall:
	$(GOCMD) install $(lib)

