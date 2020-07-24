package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//ServiceResponse is the response object for service requests
type ServiceResponse struct {
	MainURL       string
	Service       string
	AccesibleURLs []string
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG\n"))
}

//TODO: proper way of implementing random port
func availableHandler(w http.ResponseWriter, r *http.Request) {
	//Seed random
	rand.Seed(time.Now().UnixNano())

	var v int
	for i := 0; i <= 5; i++ {
		v = rand.Intn(500) + 20000
		log.Println("[AVAIL] ", v)
		l, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(v))
		if err != nil && strings.Contains(err.Error(), "address already in use") {
			continue
		}
		defer l.Close()
		break
	}

	w.Write([]byte("127.0.0.1:" + strconv.Itoa(v)))
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 Method not allowed\n"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}

	form := r.PostForm
	var service, username string

	if u, ok := form["username"]; !ok || len(u) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request: improper username"))
	} else {
		username = u[0]
	}

	// register service
	if r.Method == http.MethodPost {
		if s, ok := form["service"]; !ok || s[0] == "" {
			//use generated service name if not present on request
			service = generateServiceName()
		} else {
			service = s[0]
		}

		serviceName := service + "-" + username
		var localAddr string
		if l, ok := form["localAddr"]; !ok || len(l) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request: improper localAddr"))
		} else {
			localAddr = l[0]
		}
		log.Printf("[REGISTER] username=%s, service=%s, localAddr=%s", username, service, localAddr)
		serviceURL, err := createPortedService(&ctx, red, serviceName, "http://"+localAddr)
		if err != nil {
			//TODO: handle error
		}
		json.NewEncoder(w).Encode(&ServiceResponse{
			MainURL: serviceURL,
			Service: service,
		})
	}

	// keepalive service
	if r.Method == http.MethodPatch {

		serviceName := service + "-" + username
		if s, ok := form["service"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 Bad Request: improper service"))
		} else {
			service = s[0]
		}

		log.Printf("[ALIVE] username=%s, service=%s", username, service)
		updateKeepAlive(&ctx, red, serviceName)
	}
}
