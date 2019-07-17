package handler

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hhkbp2/go-logging"
)

type WebsocketHandler struct {
	*logging.BaseHandler
	Client *websocket.Conn
	Lock   *sync.Mutex
}

func (h *WebsocketHandler) write(messageType int, data []byte) error {
	h.Lock.Lock()
	defer h.Lock.Unlock()
	return h.Client.WriteMessage(messageType, data)
}

func (h *WebsocketHandler) Emit(record *logging.LogRecord) error {
	l := h.Format(record)
	// fmt.Println(l)
	h.write(websocket.TextMessage, []byte(l))
	return nil
}

func (h *WebsocketHandler) Handle(record *logging.LogRecord) int {
	return h.Handle2(h, record)
}

func (h *WebsocketHandler) Close() {
	h.Client.Close()
}

func NewWebsocketHandler(
	name string, addr string,
	path string, level logging.LogLevelType,
	timeout time.Duration) *WebsocketHandler {

	u := url.URL{Scheme: "ws", Host: addr, Path: path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	rv := &WebsocketHandler{
		BaseHandler: logging.NewBaseHandler(name, level),
		Client:      c,
		Lock:        &sync.Mutex{},
	}
	rv.keepalive(timeout)
	return rv
}

func (h *WebsocketHandler) keepalive(timeout time.Duration) {

	if timeout <= 0 {
		return
	}

	log.Printf("enable keepalive for %s", timeout)

	lastResponse := time.Now()
	h.Client.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			err := h.write(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				log.Print(err)
				return
			}
			time.Sleep(timeout / 2)
			diff := time.Now().Sub(lastResponse)
			if diff > timeout {
				log.Printf("diff: %s, timeout: %s, close connection", diff, timeout)
				h.Client.Close()
				return
			}
		}
	}()
}
