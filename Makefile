NAME=$(lastword $(subst /, ,$(abspath .)))
VERSION=$(shell git.exe describe --tags 2>nul || echo noversion)
GOOPT=-ldflags "-s -w -X main.version=$(VERSION)"
GOEXE=$(shell go env GOEXE)
GOOS=$(shell go env GOOS)

ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    SET=SET
else
    SET=export
endif

all:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)

_package:
	go fmt
	$(SET) "CGO_ENABLED=0" && go build $(GOOPT)
	zip $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH).zip $(NAME)$(GOEXE)

package:
	$(SET) "GOARCH=386"   && $(MAKE) _package
	$(SET) "GOARCH=amd64" && $(MAKE) _package

manifest:
	go run ./mkmanifest.go *-windows-*.zip > $(NAME).json
