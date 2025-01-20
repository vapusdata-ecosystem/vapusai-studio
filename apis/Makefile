# Makefile

# Go parameters
GOCMD = go
BUFCMD = buf

.PHONY: protogen gosync api

protogen:
	rm -rf gen-go/protos
	rm -rf gen-go/vendor
	$(BUFCMD) dep update
	$(BUFCMD) lint
	$(BUFCMD) generate
gosync: protogen
	cd gen-go && $(GOCMD) mod tidy && $(GOCMD) mod vendor
	

apisync: gosync
	cd gen-go && $(GOCMD) run main.go