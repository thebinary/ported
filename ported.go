package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/thebinary/ported/iproxy"
	"github.com/thebinary/ported/tunnel"
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
	Porter        string
	ServiceName   string
	ServerAddr    string
	Username      string
	RemoteAddr    string
	LocalAddr     string
	Timeout       time.Duration
	clientConfig  *ssh.ClientConfig
	client        *ssh.Client
	tunnel        tunnel.Tunneler
	inspectorAddr string
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

func (p *Ported) getTunnel() (err error) {
	remote, _ := net.ResolveTCPAddr("tcp", p.RemoteAddr)
	local, _ := net.ResolveTCPAddr("tcp", p.inspectorAddr)
	p.tunnel, err = tunnel.NewReverseSSH(p.client, remote, local)
	return err
}

// Start Ported connection
func (p *Ported) Start() (err error) {
	// Setup local inspector proxy
	log.Printf("==> Starging local inspector proxy")
	localURL, _ := url.Parse("http://" + p.LocalAddr)
	inspectorAddr, inspector, err := iproxy.WebInspectorMux(localURL, logReqHeader, logRespHeader, logRespBody)
	if err != nil {
		log.Fatal(err.Error())
	}
	p.inspectorAddr = inspectorAddr
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

	p.getTunnel()
	//log.Printf("sucessfully listening on tunnel server '%s' at %s", p.ServerAddr, p.RemoteAddr)

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
	p.tunnel.Connect()

	p.Close()
	return
}

// Close Ported Connections
func (p *Ported) Close() {
	log.Println("==> Cleaning up and shutting down ported...")
	log.Println("===> Closing ported tunnel...")
	if p.client != nil {
		p.client.Conn.Close()
	}
	log.Println("===> Closing ported remote listeners...")
	if p.tunnel != nil {
		p.tunnel.Close()
	}
	log.Println("==> ported successfully shutdown.")
}
