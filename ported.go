package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/thebinary/ported/flow/httpflow"
	"golang.org/x/crypto/ssh"
)

//ServiceResponse is the response object for service requests
type ServiceResponse struct {
	MainURL       string
	Service       string
	AccesibleURLs []string
}

// Ported describes an instance of a Ported
type Ported struct {
	Porter       string
	ServiceName  string
	ServerAddr   string
	Username     string
	RemoteAddr   string
	LocalAddr    string
	Timeout      time.Duration
	clientConfig *ssh.ClientConfig
	client       *ssh.Client
	listenter    net.Listener
}

// NewPorted returns a NewPorted Config object
// TODO: sanitization and validation of addresses
func NewPorted(porter, serviceName, serverAddr, username, keyFile, remoteAddr, localAddr string) (p *Ported, err error) {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("error reading key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		log.Fatalf("error parsing key: %s", keyFile)
	}
	p = &Ported{
		Porter:      porter,
		ServiceName: serviceName,
		ServerAddr:  serverAddr,
		Username:    username,
		RemoteAddr:  remoteAddr,
		LocalAddr:   localAddr,
		Timeout:     time.Second * 15,
	}
	p.clientConfig = &ssh.ClientConfig{
		User: p.Username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         p.Timeout,
	}
	return p, nil
}

func (p *Ported) keepAliveRequest() {
	formData := url.Values{
		"username":  []string{p.Username},
		"service":   []string{p.ServiceName}, // for porter server localAddr is client's remoteAddr
		"localAddr": []string{p.RemoteAddr},
	}
	http.PostForm(p.Porter+"/v1/service", formData)
}

func (p *Ported) keepAlive() {
	//TODO: get timeout from porter server
	c := time.Tick(time.Second * 40)
	for {
		<-c
		p.keepAliveRequest()
	}
}

func (p *Ported) getRemoteAddr() {
	log.Println("getting remote available address...")
	resp, err := http.Get(p.Porter + "/v1/available")
	if err != nil {
		log.Fatalf("cannot get remote available address: %v", err)
	}
	addr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("cannot get remote available address: %v", err)
	}
	p.RemoteAddr = string(addr)
}

// Start Ported connection
func (p *Ported) Start() (err error) {
	// Setup local inspector proxy
	log.Printf("==> Starging local inspector proxy")
	localURL, _ := url.Parse("http://" + p.LocalAddr)
	proxy := httputil.NewSingleHostReverseProxy(localURL)

	msg := make(chan httpflow.HTTPFlow)
	inspectTransport := DefaultInspectTransport
	inspectTransport.RequestHeaders = logReqHeader
	inspectTransport.ResponseHeaders = logRespHeader
	inspectTransport.ResponseBody = logRespBody
	inspectTransport.WebChannel = msg
	proxy.Transport = inspectTransport

	//TODO: proper way to use dynamic port
	tmpsrv := httptest.NewServer(http.NotFoundHandler())
	dynURL, _ := url.Parse(tmpsrv.URL)
	inspectorAddr := dynURL.Host
	time.Sleep(time.Second)
	tmpsrv.Close()
	inspector := NewWebInspector(msg, proxy)
	go http.ListenAndServe(inspectorAddr, inspector)

	p.getRemoteAddr()
	log.Printf("==> Starting Tunnel: %s|%s -> %s|%s", p.ServerAddr, p.RemoteAddr, inspectorAddr, p.LocalAddr)
	client, err := ssh.Dial("tcp", p.ServerAddr, p.clientConfig)
	if err != nil {
		nerr := fmt.Errorf("tunnel connect error: %v", err)
		log.Println(nerr)
		return nerr
	}
	p.client = client

	listener, err := client.Listen("tcp", p.RemoteAddr)
	if err != nil {
		client.Close()
		nerr := fmt.Errorf("error listening on '%s' on tunnel server: %v", p.RemoteAddr, err)
		log.Println(nerr)
		return nerr
	}
	p.listenter = listener
	log.Printf("sucessfully listening on tunnel server '%s' at %s", p.ServerAddr, p.RemoteAddr)

	// Call porter to create service
	formData := url.Values{
		"username":  []string{p.Username},
		"localAddr": []string{p.RemoteAddr}, // for porter server localAddr is client's remoteAddr
	}
	if serviceName != "" {
		formData.Add("service", serviceName)
	}
	resp, err := http.PostForm(p.Porter+"/v1/service", formData)
	if err != nil {
		nerr := fmt.Errorf("error requesting service: %v", err)
		log.Println(nerr)
		return nerr
	}
	service := &ServiceResponse{}
	json.NewDecoder(resp.Body).Decode(service)
	fmt.Println("\n\nYour Service is available at:\n" + service.MainURL)
	p.ServiceName = service.Service
	fmt.Printf("\nAnd the web inspector is available at:\n%s\n", "http://"+inspectorAddr+"/porter/inspector")
	// keep running keepAlive request in background
	go p.keepAlive()
	fmt.Printf("\n\n======= Logs will appear below ========\n")

	// start communication loop
	for {
		local, err := net.Dial("tcp", inspectorAddr)
		if err != nil {
			listener.Close()
			nerr := fmt.Errorf("error connecting to local address '%s': %v", p.LocalAddr, err)
			log.Println(nerr)
			return nerr
		}

		client, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting client")
		}
		//log.Printf("client connection: %s|%s -> %s", p.ServerAddr, client.RemoteAddr().String(), client.LocalAddr().String())
		handleClient(client, local)
	}
}

// Close Ported Connections
func (p *Ported) Close() {
	log.Println("==> Cleaning up and shutting down ported...")
	log.Println("===> Closing ported tunnel...")
	if p.client != nil {
		p.client.Conn.Close()
	}
	log.Println("===> Closing ported remote listeners...")
	if p.listenter != nil {
		p.listenter.Close()
	}
	log.Println("==> ported successfully shutdown.")
}
