GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOVULNCHECK = $(GOCMD) run golang.org/x/vuln/cmd/govulncheck@latest
BINARY_NAME = cocos-playable
MAIN = ./cmd/cocos-playable

.PHONY: build test race vet vulncheck verify clean

build:
	$(GOBUILD) -o ./out/$(BINARY_NAME) $(MAIN)

test:
	$(GOTEST) ./...

race:
	$(GOTEST) -race ./...

vet:
	$(GOCMD) vet ./...

vulncheck:
	$(GOVULNCHECK) ./...

verify: build test race vet vulncheck

clean:
	$(GOCMD) clean
