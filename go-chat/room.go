package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LAKuEN/go-programming-blueprint-ja/go-chat/trace"
	"github.com/gorilla/websocket"
	"github.com/stretchr/objx"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool // 操作する際にはjoin, leaveを介して操作
	tracer  trace.Tracer
}

func (r *room) run() {
	for {
		// NOTE case節のコードが同時に実行されることはない
		// これによりr.clientsへの同時操作が起こるのを防いでいる
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New user has joined.")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("A user has left.")
		case msg := <-r.forward:
			r.tracer.Trace(fmt.Sprintf("A message received: %v",
				msg.Message))
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.tracer.Trace("-> A message has been sent.")
				default:
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(`-> Failed to send the message.
					Clean up the client.`)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: messageBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatalf("Error occured at ServeHTTP: %#v", err)
	}

	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatalf("Failed to get cookie: %v", err)
		return
	}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()

	// write()とread()をそれぞれ別々に無限ループ
	go client.write()
	// メインgoroutineでread()を回し続けることで接続を保持
	client.read()
}

// 新たなroomを生成して返却
func newRoom(logging bool) *room {
	var tracer trace.Tracer
	if logging {
		tracer = trace.New(os.Stdout)
	} else {
		tracer = trace.Off()
	}
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  tracer,
	}
}
