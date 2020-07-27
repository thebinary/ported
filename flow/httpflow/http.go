package httpflow

import (
	"fmt"
	"net/http"
	"time"
)

//HTTPFlow describe the http requests and responses
type HTTPFlow struct {
	Timestamp   int64  `json:"t"`
	RemoteIP    string `json:"ip"`
	UserAgent   string `json:"u"`
	HTTPVersion string `json:"v"`
	Method      string `json:"m"`
	Host        string `json:"h"`
	Path        string `json:"p"`
	Referer     string `json:"r"`
	// Request Vars
	RequestQuery   string      `json:"rq"`
	RequestHeaders HTTPHeaders `json:"rh"`
	RequestBody    string      `json:"rb,omitifempty"`
	// Response Vars
	ResponseHeaders       HTTPHeaders `json:"rph"`
	ResponseContentType   string      `json:"rpt,omitifempty"`
	ResponseContentLength int64       `json:"rpl"`
	ResponseBody          string      `json:"rpb,omitifempty"`
	StatusCode            int         `json:"c"`
	Status                string      `json:"s"`
}

// NewHTTPFlow extracts HTTP flow info from http request and response object
func NewHTTPFlow(request *http.Request, response *http.Response) (flow *HTTPFlow) {
	//Get Remote originating IP
	var remoteIP string
	if r := request.Header.Get("X-Real-Ip"); r != "" {
		remoteIP = r
	} else {
		remoteIP = request.RemoteAddr
	}

	// Extract request headers
	rHeaders := HTTPHeaders{}
	for k, v := range request.Header {
		rHeaders[k] = v
	}

	//Extract response headers
	respHeaders := HTTPHeaders{}
	for k, v := range response.Header {
		respHeaders[k] = v
	}

	// instantiate the flow object
	flow = &HTTPFlow{
		Timestamp:             time.Now().UnixNano(),
		RemoteIP:              remoteIP,
		UserAgent:             request.UserAgent(),
		HTTPVersion:           fmt.Sprintf("%d.%d", request.ProtoMajor, request.ProtoMinor),
		Method:                request.Method,
		Path:                  request.URL.Path,
		Host:                  request.Header.Get("Host"),
		Referer:               request.Referer(),
		RequestQuery:          request.URL.RawQuery,
		RequestHeaders:        rHeaders,
		ResponseHeaders:       respHeaders,
		ResponseContentLength: response.ContentLength,
		ResponseContentType:   string(response.Header.Get("Content-Type")),
		StatusCode:            response.StatusCode,
		Status:                response.Status,
	}

	return flow
}
