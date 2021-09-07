GOCMD=go
GOBUILD=$(GOCMD) build
GOTARGET=./cmd/kubectl-multiforward.go
WINDOWS=GOOS=windows GOARCH=amd64 $(GOBUILD)
MACARM=GOOS=darwin GOARCH=arm64 $(GOBUILD)
MACAMD=GOOS=darwin GOARCH=amd64 $(GOBUILD)
LINUX=GOOS=linux GOARCH=amd64 $(GOBUILD)
GOBUILD=GOOS=windows GOARCH=amd64 $(GOBUILD)
WINDOWS_PATH=./artifacts/windows
WINDOWS_ARTIFACT=$(WINDOWS_PATH)/kubectl-multiforward
MACAMD_PATH=./artifacts/macos/amd64
MACAMD_ARTIFACT=$(MACAMD_PATH)/kubectl-multiforward
MACARM_PATH=./artifacts/macos/arm64
MACARM_ARTIFACT=$(MACARM_PATH)/kubectl-multiforward
LINUX_PATH=./artifacts/linux
LINUX_ARTIFACT=$(LINUX_PATH)/kubectl-multiforward
GOBUILD=$(GOCMD) build
GOBUILD=$(GOCMD) build
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=kubectl-multiforward

all: build

package: windows darwin-arm darwin-amd linux

windows:
	$(WINDOWS) -o $(WINDOWS_PATH)/ $(GOTARGET)
	@zip $(WINDOWS_ARTIFACT)-windows-amd64.zip $(WINDOWS_ARTIFACT).exe
	@openssl dgst -r -sha256 $(WINDOWS_ARTIFACT).exe |cut -d ' ' -f1 > $(WINDOWS_PATH)/sig
	@rm $(WINDOWS_ARTIFACT).exe

darwin-arm:
	$(MACARM) -o $(MACARM_PATH)/ $(GOTARGET)
	@tar -czvf $(MACARM_ARTIFACT)-darwin-arm64.tar.gz $(MACARM_ARTIFACT)
	@openssl dgst -r -sha256 $(MACARM_ARTIFACT) |cut -d ' ' -f1 > $(MACARM_PATH)/sig
	@rm $(MACARM_ARTIFACT)

darwin-amd:
	$(MACAMD) -o $(MACAMD_PATH)/ $(GOTARGET)
	@tar -czvf $(MACAMD_ARTIFACT)-darwin-amd64.tar.gz $(MACAMD_ARTIFACT)
	@openssl dgst -r -sha256 $(MACAMD_ARTIFACT) |cut -d ' ' -f1 > $(MACAMD_PATH)/sig
	@rm $(MACAMD_ARTIFACT)

linux:
	$(LINUX) -o $(LINUX_PATH)/ $(GOTARGET)
	@tar -czvf $(LINUX_ARTIFACT)-linux-amd64.tar.gz $(LINUX_ARTIFACT)
	@openssl dgst -r -sha256 $(LINUX_ARTIFACT) |cut -d ' ' -f1 > $(LINUX_PATH)/sig
	@rm $(LINUX_ARTIFACT)

build:
	@$(GOBUILD) -o $(BINARY_NAME) $(GOTARGET)

clean:
	@$(GOCLEAN)
	@rm -f $(BINARY_NAME)

cover:
	@$(GOTEST) -coverprofile=coverage.out
	@$(GOCMD) tool cover -html=coverage.out

.PHONY: all build clean cover linux darwin-amd darwin-arm windows package
