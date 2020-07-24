package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path"
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

func init() {
	// get current system user info
	cuUser, err := user.Current()
	if err != nil {
		log.Fatalf("error fetching current system user info: %v", err)
	}
	defKeyFile := path.Join(cuUser.HomeDir, "/.ssh/id_rsa")

	flag.StringVar(&porter, "using", "http://localhost:8888", "porter server address; API endpoint address")
	flag.StringVar(&serviceName, "as", "myapp", "service name")
	flag.StringVar(&tunnelAddr, "via", "", "porter tunnel endpoint adddress; SSH address; defaults to 'using_host:22'")
	flag.StringVar(&localAddr, "to", "127.0.0.1:8080", "local listener adddress; local service listening address")
	flag.StringVar(&userName, "with-user", cuUser.Username, "username to use for porter")
	flag.StringVar(&keyFile, "with-key", defKeyFile, "private key to use for porter")
	//TODO: this remoteAddr should be handled by porter
	flag.StringVar(&remoteAddr, "at", "127.0.0.1:8080", "listener addr at the porter tunnel endpoint")

	flag.Parse()

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
