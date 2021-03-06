package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path"
	"strconv"
	"strings"
)

var porter string
var as string
var tunnelAddr string
var localAddr string
var userName string
var keyFile string
var remoteAddr string
var serviceName string
var logReqHeader bool
var logRespHeader bool
var logRespBody bool

const flagToDescr = `local service adddress
supported address formats:
- [port]
- :[port]
- [ip]:[port]
Note: [port] format (without colon) signifies 127.0.0.1:8080
`

func sanitizeLocalAddress(addr string) (serviceAddr string, err error) {
	// if address has only digits consider it as port
	if port, err := strconv.Atoi(addr); err == nil {
		return fmt.Sprintf("127.0.0.1:%d", port), nil
	}
	_, _, err = net.SplitHostPort(addr)
	return addr, err
}

func init() {
	// get current system user info
	cuUser, err := user.Current()
	if err != nil {
		log.Fatalf("error fetching current system user info: %v", err)
	}
	defKeyFile := path.Join(cuUser.HomeDir, "/.ssh/id_rsa")

	flag.StringVar(&porter, "using", "http://localhost:8888", "porter server address; API endpoint address")
	flag.StringVar(&serviceName, "as", "", "service name; defaults to what is given by porter server")
	flag.StringVar(&tunnelAddr, "via", "", "porter tunnel endpoint adddress; SSH address; defaults to 'using_host:22'")
	flag.StringVar(&localAddr, "to", "127.0.0.1:8080", flagToDescr)
	flag.StringVar(&userName, "with-user", cuUser.Username, "username to use for porter")
	flag.StringVar(&keyFile, "with-key", defKeyFile, "private key to use for porter")
	//TODO: this remoteAddr should be handled by porter
	flag.StringVar(&remoteAddr, "at", "127.0.0.1:8080", "listener addr at the porter tunnel endpoint")
	flag.BoolVar(&logReqHeader, "log.request.header", false, "log request headers")
	flag.BoolVar(&logRespHeader, "log.response.header", false, "log response headers")
	flag.BoolVar(&logRespBody, "log.response.body", false, "log request body")
	flag.Parse()

	localAddr, err = sanitizeLocalAddress(localAddr)
	if err != nil {
		fmt.Errorf("invalid form of -to address: ")
		os.Exit(1)
	}

	// extract default tunnelAddr from porter url
	// eg: http://ported.example.com:8888 -> extract and convert -> ported.example.com:22
	if tunnelAddr == "" {
		httpParts := strings.Split(porter, "//")
		addrs := strings.Split(httpParts[1], ":")
		tunnelAddr = addrs[0] + ":22"
	}
	//log.Printf("[SERVICE] porter=%s, username=%s, service=%s, localAddr=%s, tunnelAddr=%s", porter, userName, serviceName, localAddr, tunnelAddr)
}

func main() {
	//os.Exit(2)
	p, err := NewPorted(porter, serviceName, tunnelAddr, userName, keyFile, remoteAddr, localAddr)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go p.Start()

	// block until stop signal is recieved
	<-c
	p.Close()
}
