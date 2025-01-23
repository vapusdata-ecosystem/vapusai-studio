# Makefile

#input params
OCI_IMAGE_PATH = asia-south1-docker.pkg.dev/vapusdata-beta/vapusai-studio/aitools
TAG = local_$(shell whoami)

# Go parameters
BUFCMD = buf
POETRYCMD = poetry
BANDITCMD = bandit
PYLINTCMD = pylint
PYINSTALLER = pyinstaller

BINARY_NAME = aitools

all: test build

protogen:
	cd ../../../apis && $(BUFCMD) generate

.PHONY: all build test clean run

build-binary: protogen 
	$(POETRYCMD) run $(PYINSTALLER) server.spec --clean --noconfirm

build: build-binary
	docker build --no-cache -t $(OCI_IMAGE_PATH):$(TAG) .
	rm -f $(BINARY_NAME)

release: build
	docker push $(OCI_IMAGE_PATH):$(TAG) 

lint:
	$(POETRYCMD) run $(PYLINTCMD) controller/ models/ services/ utils/ core/ datastores/ helpers/ --rcfile=pylint.rc

sec:
	$(POETRYCMD) run $(BANDITCMD) -c bandit.toml -r -f json .

prehook: lint sec

runmain:
	poetry run python server.py --conf=$(CONF) --debug=true