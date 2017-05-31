package main

import (
    "bytes"
    "flag"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait = 10 * time.Second

    pongWait = 60 * time.Second

    pingPeriod = (pongWait * 9) / 10

    maxMessageSize = 1024 * 1024
)


var (
    newline = []byte{'\n'}
    space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  maxMessageSize,
    WriteBufferSize: maxMessageSize,
}

type Client struct {
    chatServer *ChatServer

    ws *websocket.Conn

    send chan []byte

}

func (client *Client) readPump() {
    defer func() {
        client.chatServer.unregister <- client
        client.ws.Close()
    }()
    client.ws.SetReadLimit(maxMessageSize)
    client.ws.SetReadDeadline(time.Now().Add(pongWait))
    client.ws.SetPongHandler(func(string) error { client.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        _, message, err := client.ws.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
                log.Printf("error: %v", err)
            }
            break
        }
        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
        client.chatServer.broadcast <- message
    }
}

func (client *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        client.ws.Close()
    }()
    for {
        select {
        case message, ok := <-client.send:
            client.ws.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {

                client.ws.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := client.ws.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            n := len(client.send)
            for i := 0; i < n; i++ {
                w.Write(newline)
                w.Write(<-client.send)
            }

            if err := w.Close(); err != nil {
                return
            }
        case <-ticker.C:
            client.ws.SetWriteDeadline(time.Now().Add(writeWait))
            if err := client.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
                return
            }
        }
    }
}

func serveWs(chatServer *ChatServer, writer http.ResponseWriter, request *http.Request) {
    ws, err := upgrader.Upgrade(writer, request, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{chatServer: chatServer, ws: ws, send: make(chan []byte, 256)}
    client.chatServer.register <- client
    go client.writePump()
    client.readPump()
}

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(writer http.ResponseWriter, request *http.Request) {
    log.Println(request.URL)
    http.ServeFile(writer, request, "html/chat.html")
}

func chat() {
    flag.Parse()
    chatServer := newChatServer()
    go chatServer.run()
    http.HandleFunc("/", serveHome)
    http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
        serveWs(chatServer, writer, request)
    })
    err := http.ListenAndServe(*addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

type ChatServer struct {
    clients map[*Client]bool

    broadcast chan []byte

    register chan *Client

    unregister chan *Client
}

func newChatServer() *ChatServer {
    return &ChatServer{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

func (chatserver *ChatServer) run() {
    for {
        select {
        case client := <-chatserver.register:
            chatserver.clients[client] = true
        case client := <-chatserver.unregister:
            if _, ok := chatserver.clients[client]; ok {
                delete(chatserver.clients, client)
                close(client.send)
            }
        case message := <-chatserver.broadcast:
            for client := range chatserver.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(chatserver.clients, client)
                }
            }
        }
    }
}