package handler

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hhkbp2/go-logging"
)

type WebsocketHandler struct {
	*logging.BaseHandler
	Client *websocket.Conn
}

func (h *WebsocketHandler) Emit(record *logging.LogRecord) error {
	l := h.Format(record)
	// fmt.Println(l)
	h.Client.WriteMessage(websocket.TextMessage, []byte(l))
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
	keepaliveInterval int) *WebsocketHandler {

	u := url.URL{Scheme: "ws", Host: addr, Path: path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// keepalive
	if keepaliveInterval > 0 {
		keepalive(c, time.Second*time.Duration(keepaliveInterval))
	}

	return &WebsocketHandler{
		BaseHandler: logging.NewBaseHandler(name, level),
		Client:      c,
	}
}

func keepalive(c *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				log.Print(err)
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}
