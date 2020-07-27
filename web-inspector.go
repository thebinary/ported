//go:generate go-bindata -fs -prefix "inspector/dist/" -o web-inspector-fs.go inspector/dist/...
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/thebinary/ported/flow/httpflow"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func WebSocketHandlerWithMessageChannel(msg chan httpflow.HTTPFlow) func(http.ResponseWriter, *http.Request) {
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

func NewWebInspector(ch chan httpflow.HTTPFlow, proxy *httputil.ReverseProxy) (inspectorWeb *http.ServeMux) {
	inspectorWeb = http.NewServeMux()
	inspectorWeb.Handle("/", proxy)
	inspectorWeb.HandleFunc("/porter/stream", WebSocketHandlerWithMessageChannel(ch))
	webDir := AssetFile()
	inspectorWeb.Handle("/porter/inspector/", http.StripPrefix("/porter/inspector", http.FileServer(webDir)))
	return inspectorWeb
}
