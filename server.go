package top_senders

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/textileio/go-threads/broadcast"
)

var rootDir fs.FS

func Serve(port string) {
	rootDir, _ = fs.Sub(content, "static")
	latestGraph := []byte{'{', '}'}

	var cast broadcast.Broadcaster
	var upgrader = websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	upgrader.EnableCompression = true

	go func() {
		for update := range Updates {
			data, err := json.Marshal(update)
			if err != nil {
				log.Println(err)
				continue
			}
			latestGraph = data
			_ = cast.Send(latestGraph)
		}
	}()

	http.HandleFunc("/latest", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(latestGraph)
	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		c, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			return
		}
		defer c.Close()
		sub := cast.Listen()
		defer sub.Discard()
		for message := range sub.Channel() {
			e := c.WriteMessage(websocket.TextMessage, message.([]byte))
			if e != nil {
				return
			}
		}
	})

	http.Handle("/", &CacheHandler{})
	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()

	cast.Discard()
	log.Fatal(err)
}

// CacheHandler implements the Handler interface with a Cache-Control set on responses
type CacheHandler struct{}

func (ch CacheHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Cache-Control", "public, max-age=3600")
	http.FileServer(http.FS(rootDir)).ServeHTTP(writer, request)
}
