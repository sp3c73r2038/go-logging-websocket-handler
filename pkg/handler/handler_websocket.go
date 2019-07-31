package handler

import (
	"log"
	"net/url"
	"sync"

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
	name string, u url.URL, level logging.LogLevelType) *WebsocketHandler {

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	rv := &WebsocketHandler{
		BaseHandler: logging.NewBaseHandler(name, level),
		Client:      c,
		Lock:        &sync.Mutex{},
	}
	return rv
}
