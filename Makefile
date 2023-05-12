include .env

PROJECTNAME=$(shell basename "$(PWD)")

# Go переменные.
GOBASE=$(shell pwd)
GOPATH=$(GOBASE)/vendor:$(GOBASE)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)
#SERVER=$()

go-build:
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBIN)/$(PROJECTNAME)

go-deploy:
	@echo "  >  Building binary..."
	GOOS=linux GOARCH=amd64 go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)
	$(shell ssh root@$(SERVER) "ps aux")# | grep $(PROJECTNAME) | awk '{print $2}' | xargs kill -9")
	#$(shell scp $(GOBIN)/$(PROJECTNAME) root@$(SERVER):/root/$(PROJECTNAME))