package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var channels = make(map[string]string)
var clients = make(map[string]User)

type User struct {
	DisplayName string
	Connection  *websocket.Conn
}

type HTMLDir struct {
	d http.Dir
}

func (d HTMLDir) Open(name string) (http.File, error) {
	// Routes without file extension are HTML
	f, err := d.d.Open(name + ".html")

	// If the .html extension of name doesn't exist
	// we try opening as is
	if os.IsNotExist(err) {
		if f, err := d.d.Open(name); err == nil {
			return f, nil
		}
	}

	return f, err
}

// Upgrade and manage the websocket connection
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	// If connection could not be upgraded
	if err != nil {
		fmt.Println("here")
		return
	}

	// While connection is active
	for {
		// Open incoming message
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// do something with message here

		// Send message to browser, or return if error
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

func main() {
	// Router for static Svelte frontend
	svelte := http.FileServer(HTMLDir{http.Dir("./client/build/")})
	http.Handle("/", http.StripPrefix("/", svelte))

	// Websocket connection
	http.HandleFunc("/ws", Websocket)

	http.ListenAndServe(":8080", nil)
}
