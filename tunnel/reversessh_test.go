package tunnel

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"strconv"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

var username string
var sshAddr *net.TCPAddr
var remote *net.TCPAddr
var local *net.TCPAddr
var sshConfig *ssh.ClientConfig

func currentUserAndKey() (username, key string, err error) {
	cuUser, err := user.Current()
	username = cuUser.Username
	if err != nil {
		return username, key, err
	}
	key = path.Join(cuUser.HomeDir, ".ssh", "id_rsa")
	return username, key, nil
}

func parseTCPAddr(addr string) (tcpAddr *net.TCPAddr, err error) {
	if host, port, err := net.SplitHostPort(addr); err != nil {
		return nil, err
	} else {
		portInt, _ := strconv.Atoi(port)
		parsedIP := net.ParseIP(host)
		if parsedIP == nil {
			return nil, err
		}
		tcpAddr = &net.TCPAddr{
			IP:   parsedIP,
			Port: portInt,
		}
	}
	return tcpAddr, nil
}

func init() {
	var key string
	var err error
	var ok bool

	// Set SSH Address for tests
	if addr, ok := os.LookupEnv("TEST_SSHADDR"); ok {
		sshAddr, err = parseTCPAddr(addr)
		if err != nil {
			log.Fatalf("error parsing SSH addr")
		}
	} else {
		sshAddr = &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 22,
		}
	}
	// Set Remote Tunnel Address for tests
	if addr, ok := os.LookupEnv("TEST_REMOTE"); ok {
		remote, err = parseTCPAddr(addr)
		if err != nil {
			log.Fatalf("error parsing Remote addr")
		}
	} else {
		remote = &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 7777,
		}
	}

	// Set Local Tunnel Address for tests
	if addr, ok := os.LookupEnv("TEST_LOCAL"); ok {
		local, err = parseTCPAddr(addr)
		if err != nil {
			log.Fatalf("error parsing Local addr")
		}
	} else {
		local = &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: 8080,
		}
	}

	// Set username for tests
	if username, ok = os.LookupEnv("TEST_USERNAME"); !ok {
		username, key, err = currentUserAndKey()
		if err != nil {
			log.Fatal("no TEST_USERNAME suplied to be used for SSH, also could not get current system username")
		}
	}

	// parse SSH private key
	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		log.Fatalf("error reading key '%s': %v", key, err)
	}
	sshKey, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		log.Fatalf("error parsing SSH key '%s': %v", key, err)
	}

	sshConfig = &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sshKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
}

func TestReverseSSH(t *testing.T) {
	client, err := ssh.Dial(sshAddr.Network(), sshAddr.String(), sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	r, err := NewReverseSSH(
		client,
		remote,
		local,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(r.remoteEndpoint, r.localEndpoint)

	err = r.Connect()
	if err != nil {
		log.Fatalf(err.Error())
	}
	t.Logf("tunnel: %+v", r)
}
