package handler

import (
	"log"
	"net/url"

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
	path string, level logging.LogLevelType) *WebsocketHandler {

	u := url.URL{Scheme: "ws", Host: addr, Path: path}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return &WebsocketHandler{
		BaseHandler: logging.NewBaseHandler(name, level),
		Client:      c,
	}
}
