NAME=mist_exporter
VERSION=$(shell cat VERSION)

PWD=$(shell pwd)
BINDIR=${PWD}/bin
CMDDIR=${PWD}/cmd


.PHONY: ci
ci: vet test lint


.PHONY: clean
clean:
	go clean
	rm -rf ${BINDIR}/*


.PHONY: dep
dep:
	go get -u ./... && go mod tidy


.PHONY: build
build: dep vet test lint
	go build  -C ${CMDDIR}/ -o ${BINDIR}/${NAME}-${VERSION}


.PHONY: run
run: build
	${BINDIR}/${NAME}-${VERSION}


.PHONY: test
test:
	go test -v ./... -count=1 --race


.PHONY: test_coverage
test_coverage:
	go test ./... -coverprofile=coverage.out


.PHONY: vet
vet:
	go vet ./...


.PHONY: lint
lint:
	golangci-lint run --enable-all
