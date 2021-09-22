TESTABLE_PACKAGES = `go list ./... | grep -v examples | grep -v constants | grep -v mocks | grep -v helpers | grep -v interfaces | grep -v protos | grep -v e2e | grep -v benchmark`
GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
BINARY_NAME = Cocos3DPlayableAdPlugin
setup: init-submodules
	@go get ./...

init-submodules:
	@git submodule init

setup-ci:
	@go get github.com/mattn/goveralls
	@go get -u github.com/wadey/gocovmerge

setup-protobuf-macos:
	@brew install protobuf
	@go get github.com/golang/protobuf/protoc-gen-go

run:
	@go run main.go

runCli:
	@cd ./bin && .\go_build_c3d2one