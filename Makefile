ROOTDIR=/opt/rocinax
APPNAME=rhole
APPVERSION=0.0.0-alpha

# Go build commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Develop Environment Directory
CMDDIR=cmd
BUILDWORKDIR=work
RELEASEDIR=release

# Production Environment Directory
BINARYDIR=bin
TEMPLATEDIR=template
CONFIGDIR=config
LOGDIR=logs
INITDIR=init

$(eval CMDS := $(shell find $(CMDDIR)/* -type d | sed 's/$(CMDDIR)\///'))

.PHONY: all
all: test release

# make test ... test product
.PHONY: test
test:
	$(GOTEST) -v ./...

# make install ... install application
.PHONY: install
install: clean build
	mkdir -p $(ROOTDIR)/$(APPNAME)
	mkdir -p $(ROOTDIR)/$(APPNAME)
	cp -Rp $(BUILDWORKDIR)/$(APPNAME) /opt/$(APPNAME)/

# make clean ... clean application directory
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILDWORKDIR)/*

# make build ... build application
.PHONY: build
build:
	mkdir -p $(BUILDWORKDIR)/$(APPNAME)/$(BINARYDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPNAME)/$(TEMPLATEDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPNAME)/$(CONFIGDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPNAME)/$(LOGDIR)
	mkdir -p $(BUILDWORKDIR)/$(APPNAME)/$(INITDIR)
	cp -Rp $(TEMPLATEDIR) $(BUILDWORKDIR)/$(APPNAME)/
	cp -Rp $(CONFIGDIR) $(BUILDWORKDIR)/$(APPNAME)/
	cp -Rp $(INITDIR) $(BUILDWORKDIR)/$(APPNAME)/
	$(foreach CMD,$(CMDS),CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILDWORKDIR)/$(APPNAME)/$(BINARYDIR)/$(CMD) -v $(CMDDIR)/$(CMD)/main.go)
	
# release source files
.PHONY: release
release: clean build
	rm -f $(RELEASEDIR)/$(APPNAME)-$(APPVERSION).tar.gz
	mkdir -p $(RELEASEDIR)
	tar -zcvf $(RELEASEDIR)/$(APPNAME)-$(APPVERSION).tar.gz -C $(BUILDWORKDIR)/ .

# cross compile
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARYUNIX) -v