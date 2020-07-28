package tunnel

//TODO: implementation for network other than tcp like unix sockets

import (
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

/*
ReverseSSH defines the reverse ssh tunnel.

This is a golang implementation for the reverse SSH tunnel
that can be created with 'ssh' client binary with:
ssh -R remote_port:localhost:local_port
*/
type ReverseSSH struct {
	client         *ssh.Client  //SSH client to use for the reverse tunnel
	remoteEndpoint net.Addr     // RemoteEndpoint is the remote endpoint of the tunnel
	localEndpoint  net.Addr     // LocalEndpoint is the local endpoint of the tunnel
	listener       net.Listener // Listener at the remote end of the tunnel which will forward the requests
}

//Client is the getter for the underlying ssh client
func (r *ReverseSSH) Client() *ssh.Client {
	return r.client
}

//Listen starts the remote endpoint listener
func (r *ReverseSSH) listen() (err error) {
	r.listener, err = r.client.Listen(r.remoteEndpoint.Network(), r.remoteEndpoint.String())
	return err
}

// NewReverseSSH returns a reverse SSH tunnel instance
func NewReverseSSH(client *ssh.Client, remote, local net.Addr) (r *ReverseSSH, err error) {
	r = &ReverseSSH{
		client:         client,
		remoteEndpoint: remote,
		localEndpoint:  local,
	}
	return r, nil
}

func (r *ReverseSSH) forward(client net.Conn, remote net.Conn) {
	//defer client.Close()
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
		}
	}()

	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
		}
	}()
}

/*
Connect starts the loop for accepting connections at remote endpoint and forwarding them to
the local endpoint of the tunnel
*/
func (r *ReverseSSH) Connect() (err error) {
	if err := r.listen(); err != nil {
		return err
	}

	for {
		localConn, err := net.Dial(r.localEndpoint.Network(), r.localEndpoint.String())
		if err != nil {
			r.listener.Close()
			nerr := fmt.Errorf("error connecting to local address '%s': %v", r.localEndpoint.String(), err)
			log.Println(nerr)
			return nerr
		}
		client, err := r.listener.Accept()
		if err != nil {
			log.Printf("error accepting client")
		}
		r.forward(client, localConn)
	}
}

/*
Close closed the underlying the remote listener
Users must remember to call this to release the resources
to avoid remote listener remain as leftover.
*/
func (r *ReverseSSH) Close() (err error) {
	return r.listener.Close()
}
