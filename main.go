package main

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var chatroom = MakeChatroom()

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
	// Try to upgrade and return if we can't
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Try to register the user and return otherwise
	chatroom.Register("demo-session", "CoolDude123")
	if err != nil {
		return
	}

	// Try to connect and return if otherwise
	chatroom.Connect("demo-session", conn)
	if err != nil {
		return
	}

	// While connection is active
	for {
		// Open incoming message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		err = chatroom.ReceiveMessage("demo-session", string(msg), "main")
		if err != nil {
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
