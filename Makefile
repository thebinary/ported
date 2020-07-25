MAKE		= make
SRCS 		= main.go middleman.go ported.go
LDFLAGS 	= -s -w

all: build/ported-linux_amd64 build/ported-darwin ported_all

build/ported-linux_amd64: ${SRCS}
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o $@ ${SRCS}
	upx $@
build/ported-darwin: ${SRCS}
	GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o $@ ${SRCS}
	upx $@

ported_all:
	cd ./porter; $(MAKE) all

clean:
	@echo "Removing ported binaries..."
	rm -f build/ported-*
	@echo ""
	@echo "Removing porter binaries..."
	cd ./porter; $(MAKE) clean