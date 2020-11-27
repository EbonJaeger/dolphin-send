PKGNAME=dolphin-send
MODULE=github.com/EbonJaeger/dolphin-send
VERSION="1.0.0"

GO?=go
GOFLAGS?=
GOOS?=linux
GOARCH?=amd64

GOSRC!=find . -name '*.go'
GOSRC+=go.mod go.sum

dolphin-send: $(GOSRC)
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) \
		-ldflags " \
		-X main.Version=$(VERSION)" \
		-o $(PKGNAME) \
		./cmd/$(PKGNAME)

all: dolphin-send

RM?=rm -f

clean:
	$(GO) mod tidy
	$(RM) $(PKGNAME)
	$(RM) -r vendor

install: all
	$(GO) install ./cmd/$(PKGNAME)

uninstall:
	$(RM) $(BINDIR)/$(PKGNAME)

check:
	$(GO) get -u golang.org/x/lint/golint
	$(GO) get -u github.com/securego/gosec/cmd/gosec
	$(GO) get -u honnef.co/go/tools/cmd/staticcheck
	$(GO) get -u gitlab.com/opennota/check/cmd/aligncheck
	$(GO) fmt -x ./...
	$(GO) vet ./...
	golint -set_exit_status `go list ./... | grep -v vendor`
	gosec -exclude=G204,G307 ./...
	staticcheck ./...
	aligncheck ./...
	$(GO) test -cover ./...

vendor: check clean
	$(GO) mod vendor

.DEFAULT_GOAL := all

.PHONY: all clean install uninstall check vendor
