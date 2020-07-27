package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/thebinary/ported/flow/httpflow"
)

//InspectTransport structure
type InspectTransport struct {
	RequestHeaders  bool
	ResponseHeaders bool
	ResponseBody    bool
	WebChannel      chan httpflow.HTTPFlow
}

//DefaultInspectTransport is the default inspect transport instance initialized by init
var DefaultInspectTransport *InspectTransport

func init() {
	DefaultInspectTransport = &InspectTransport{}
}

//NewInspectTransport returns a new instance of InspectTranspor
//TODO: optimized transport
func NewInspectTransport(responseHeaders, responseBody, requestHeaders bool) (transport *InspectTransport) {
	return &InspectTransport{
		ResponseHeaders: responseHeaders,
		ResponseBody:    responseBody,
		RequestHeaders:  requestHeaders,
	}
}

//RoundTrip is the implementation method for http RoundTripper
func (i *InspectTransport) RoundTrip(request *http.Request) (response *http.Response, err error) {
	start := time.Now()
	response, err = http.DefaultTransport.RoundTrip(request)
	elapsed := time.Since(start)

	w := *httpflow.NewHTTPFlow(request, response)
	// default loggin
	// eg: 2020/07/25 18:30:26 1.1.1.1 4.046181ms   GET "/test" HTTP/1.1 200 839 "" "curl/7.54.0"
	accessLog := fmt.Sprintf("%s %-12s %s \"%s\" HTTP/%d.%d %d %d \"%s\" \"%s\"",
		w.RemoteIP, elapsed,
		request.Method, request.URL.Path, request.ProtoMajor, request.ProtoMinor,
		response.StatusCode, response.ContentLength,
		request.Referer(), request.UserAgent())
	log.Println(accessLog)

	if i.RequestHeaders {
		fmt.Println("")
		fmt.Println("---- REQUEST ----")
		fmt.Println(w.RequestHeaders)
		fmt.Println("-----------------")
		fmt.Println("")
	}

	//TODO: [FIX] handle headerOnly or bodyOnly cases for logging
	if i.ResponseHeaders || i.ResponseBody {
		fmt.Println("")
		fmt.Println("---- RESPONSE ----")
		fmt.Println(w.ResponseHeaders)
		fmt.Println("------------------")
		fmt.Println("")
	}

	if i.WebChannel != nil {
		go func() {
			i.WebChannel <- w
		}()
	}

	return response, err
}
