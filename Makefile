GO_EXECUTABLE := go
VERSION := $(shell git describe --abbrev=10 --dirty --always --tags)
DIST_DIRS := find * -type d -exec

all: build install

build:
	${GO_EXECUTABLE} build -o fastqutils -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils

install:
	${GO_EXECUTABLE} install -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils

test:
	${GO_EXECUTABLE} test github.com/fredericlemoine/fastqutils/tests/

.PHONY: deploy

deploy:
	mkdir -p deploy/${VERSION}
	env GOOS=windows GOARCH=amd64 ${GO_EXECUTABLE} build -o deploy/${VERSION}/fastqutils_amd64.exe -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils
	env GOOS=windows GOARCH=386 ${GO_EXECUTABLE} build -o deploy/${VERSION}/fastqutils_386.exe -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils
	env GOOS=darwin GOARCH=amd64 ${GO_EXECUTABLE} build -o deploy/${VERSION}/fastqutils_amd64_darwin -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils
	env GOOS=linux GOARCH=amd64 ${GO_EXECUTABLE} build -o deploy/${VERSION}/fastqutils_amd64_linux -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils
	env GOOS=linux GOARCH=386 ${GO_EXECUTABLE} build -o deploy/${VERSION}/fastqutils_386_linux -ldflags "-X github.com/fredericlemoine/fastqutils/cmd.Version=${VERSION}" github.com/fredericlemoine/fastqutils
	tar -czvf deploy/${VERSION}.tar.gz --directory="deploy" ${VERSION}
