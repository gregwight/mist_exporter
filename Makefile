NAME=mist_exporter
VERSION=0.0.1

OSES=linux
ARCHES=amd64 arm64

PWD=$(shell pwd)
BINDIR=${PWD}/bin
CMDDIR=${PWD}/cmd
DISTDIR=${PWD}/dist
GZCMD = tar -czf
ZIPCMD = zip -r


.PHONY: all
all: dist


.PHONY: clean
clean:
	go clean
	rm -rf ${BINDIR}/*
	rm -rf ${DISTDIR}/* 


.PHONY: dep
dep:
	go get -u ./... && go mod tidy


.PHONY: dist
dist: clean dep build
	for OS in ${OSES}; do \
		for ARCH in ${ARCHES}; do \
			BUILDNAME=${NAME}-${VERSION}-$${OS}-$${ARCH}; \
			GOARCH=$${ARCH} GOOS=$${OS} go build -C ${CMDDIR}/ -o ${DISTDIR}/$${BUILDNAME}/${NAME}; \
			cd ${DISTDIR}; \
			cp $(PWD)/config.yaml $${BUILDNAME}/; \
			if [ "$${OS}" = "windows" ]; then \
				${ZIPCMD} $${BUILDNAME}.zip $${BUILDNAME}; \
			else \
				${GZCMD} $${BUILDNAME}.tar.gz $${BUILDNAME}; \
			fi; \
			rm -rf $${BUILDNAME}; \
			cd ${PWD}; \
		done; \
	done


.PHONY: dep
build: dep vet test
	go build  -C ${CMDDIR}/ -o ${BINDIR}/${NAME}


.PHONY: run
run: build
	${BINDIR}/${NAME}


.PHONY: run-dev
run-dev:
	go run . -d

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
