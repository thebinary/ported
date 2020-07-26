package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

//InspectTransport structure
type InspectTransport struct {
	RequestHeaders  bool
	ResponseHeaders bool
	ResponseBody    bool
	WebChannel      chan webLog
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

	//Get Remote originating IP
	var remoteIP string
	if r := request.Header.Get("X-Real-Ip"); r != "" {
		remoteIP = r
	} else {
		remoteIP = request.RemoteAddr
	}

	// Extract request headers
	rHeaders := map[string][]string{}
	for k, v := range request.Header {
		rHeaders[k] = v
	}

	//Extract response headers
	respHeaders := map[string][]string{}
	for k, v := range response.Header {
		respHeaders[k] = v
	}

	// default loggin
	// eg: 2020/07/25 18:30:26 1.1.1.1 4.046181ms   GET "/test" HTTP/1.1 200 839 "" "curl/7.54.0"
	accessLog := fmt.Sprintf("%s %-12s %s \"%s\" HTTP/%d.%d %d %d \"%s\" \"%s\"",
		remoteIP, elapsed,
		request.Method, request.URL.Path, request.ProtoMajor, request.ProtoMinor,
		response.StatusCode, response.ContentLength,
		request.Referer(), request.UserAgent())
	log.Println(accessLog)
	w := webLog{
		Timestamp:             time.Now().Unix(),
		ResponseTime:          elapsed.String(),
		Method:                request.Method,
		Path:                  request.URL.Path,
		HTTPVersion:           fmt.Sprintf("%d.%d", request.ProtoMajor, request.ProtoMinor),
		StatusCode:            response.StatusCode,
		Status:                response.Status,
		ResponseContentLength: response.ContentLength,
		Referer:               request.Referer(),
		UserAgent:             request.UserAgent(),
		RemoteIP:              remoteIP,
		RequestHeaders:        rHeaders,
		ResponseHeaders:       respHeaders,
	}

	if req, err := httputil.DumpRequestOut(request, false); err == nil {
		reqStr := string(req)
		if i.RequestHeaders {
			fmt.Println("")
			fmt.Println("---- REQUEST ----")
			fmt.Println(reqStr)
			fmt.Println("-----------------")
			fmt.Println("")
		}
	}

	if resp, err := httputil.DumpResponse(response, true); err == nil {
		respStr := string(resp)
		//TODO: [FIX] handle headerOnly or bodyOnly cases for logging
		if i.ResponseHeaders || i.ResponseBody {
			fmt.Println("")
			fmt.Println("---- RESPONSE ----")
			fmt.Println(respStr)
			fmt.Println("------------------")
			fmt.Println("")
		}
	}

	if i.WebChannel != nil {
		go func() {
			i.WebChannel <- w
		}()
	}

	return response, err
}
