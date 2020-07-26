MAKE		= make
SRCS 		= main.go middleman.go ported.go inspector.go web-inspector.go
GEN_SRCS 	= web-inspector-fs.go
LDFLAGS 	= -s -w

all: build/ported-linux_amd64 build/ported-darwin ported_all

${GEN_SRCS}:
	go generate .

build/ported-linux_amd64: ${SRCS} ${GEN_SRCS}
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o $@ ${SRCS} ${GEN_SRCS}
	upx $@
build/ported-darwin: ${SRCS} ${GEN_SRCS}
	GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o $@ ${SRCS} ${GEN_SRCS}
	upx $@

ported_all:
	cd ./porter; $(MAKE) all

gen-clean:
	rm -f ${GEN_SRCS}

install: ${SRCS} ${GEN_SRCS}
	go install

clean:
	@echo "Removing ported binaries..."
	rm -f build/ported-*
	@echo ""
	@echo "Removing porter binaries..."
	cd ./porter; $(MAKE) clean