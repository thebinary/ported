SRCS 		= porter.go redis.go handlers.go
LDFLAGS 	= -s -w

all: ../build/porter-linux_amd64 ../build/porter-darwin
../build/porter-linux_amd64: ${SRCS}
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o $@ ${SRCS}
	upx $@
../build/porter-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o $@ ${SRCS}
	upx $@

clean:
	rm -f ../build/porter-*