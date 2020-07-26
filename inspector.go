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

//NewInspectTransport returns a new instance of InspectTransport
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

	// default loggin
	// eg: 2020/07/25 18:30:26 4.046181ms   GET "/test" HTTP/1.1 200 839 "" "curl/7.54.0"
	accessLog := fmt.Sprintf("%-12s %s \"%s\" HTTP/%d.%d %d %d \"%s\" \"%s\"", elapsed,
		request.Method, request.URL.Path, request.ProtoMajor, request.ProtoMinor,
		response.StatusCode, response.ContentLength,
		request.Referer(), request.UserAgent())
	log.Println(accessLog)

	if i.WebChannel != nil {
		go func() {
			i.WebChannel <- webLog{
				Timestamp:             time.Now().Unix(),
				ResponseTime:          elapsed.String(),
				Method:                request.Method,
				Path:                  request.URL.Path,
				HTTPVersion:           fmt.Sprintf("%d.%d", request.ProtoMajor, request.ProtoMinor),
				StatusCode:            response.StatusCode,
				ResponseContentLength: response.ContentLength,
				Referer:               request.Referer(),
				UserAgent:             request.UserAgent(),
			}
		}()
	}

	if i.RequestHeaders {
		if req, err := httputil.DumpRequestOut(request, false); err == nil {
			fmt.Println("")
			fmt.Println("---- REQUEST ----")
			fmt.Println(string(req))
			fmt.Println("-----------------")
			fmt.Println("")
		}
	}

	if i.ResponseHeaders || i.ResponseBody {
		if resp, err := httputil.DumpResponse(response, i.ResponseBody); err == nil {
			fmt.Println("")
			fmt.Println("---- RESPONSE ----")
			fmt.Println(string(resp))
			fmt.Println("------------------")
			fmt.Println("")
		}
	}

	return response, err
}
