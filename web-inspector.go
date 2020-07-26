//go:generate go-bindata -fs -prefix "inspector/dist/" -o web-inspector-fs.go inspector/dist/...
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type webLog struct {
	Timestamp             int64  `json:"t"`
	ResponseTime          string `json:"rtt"`
	Method                string `json:"m"`
	Path                  string `json:"p"`
	HTTPVersion           string `json:"h"`
	StatusCode            int    `json:"c"`
	Referer               string `json:"r"`
	ResponseContentLength int64  `json:"rpl"`
	UserAgent             string `json:"u"`
}

func WebSocketHandlerWithMessageChannel(msg chan webLog) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost:8081"},
		})
		if err != nil {
			log.Println(err)
			return
		}
		defer c.Close(websocket.StatusInternalError, "internal server error")

		ctx := context.Background()
		ctx = c.CloseRead(ctx)

		for {
			select {
			case m := <-msg:
				j, _ := json.Marshal(m)
				wsjson.Write(ctx, c, string(j))
			case <-ctx.Done():
				c.Close(websocket.StatusNormalClosure, "")
				return
			}
		}
	}
}

func NewWebInspector(ch chan webLog, proxy *httputil.ReverseProxy) (inspectorWeb *http.ServeMux) {
	inspectorWeb = http.NewServeMux()
	inspectorWeb.Handle("/", proxy)
	inspectorWeb.HandleFunc("/porter/stream", WebSocketHandlerWithMessageChannel(ch))
	webDir := AssetFile()
	inspectorWeb.Handle("/porter/inspector/", http.StripPrefix("/porter/inspector", http.FileServer(webDir)))
	return inspectorWeb
}
