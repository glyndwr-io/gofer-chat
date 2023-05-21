package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
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

// Make sure requests have a session cookie and redirect
// index requests to login if not registered
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionID string

		// Get the sessionID from cookie, or
		// set if it does not exist
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sessionID = uuid.NewString()
			r.AddCookie(&http.Cookie{
				Name:  "session_id",
				Value: sessionID,
			})
		} else {
			sessionID = cookie.Value
		}

		// Can't redirect other static files
		if r.URL.Path != "/" {
			next.ServeHTTP(w, r)
			return
		}

		// Redirect if not registered
		registered := chatroom.IsRegistered(sessionID)
		fmt.Println(registered)
		if err != nil || !registered {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Attempt to register a user with the chat service
func Register(w http.ResponseWriter, r *http.Request) {
	// Even though we set the cookie, user may have
	// cookies disabled
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Cookies are disabled", http.StatusUnauthorized)
		return
	}
	sessionID := cookie.Value

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	displayName := r.FormValue("display-name")

	// Try to register the user and return otherwise
	err = chatroom.Register(sessionID, displayName)
	if err != nil {
		http.Error(w, "User already exists", http.StatusUnauthorized)
		return
	}
}

// Upgrade and manage the websocket connection
func Websocket(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	sessionID := cookie.Value

	// Try to upgrade and return if we can't
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("here1")
		return
	}

	// Try to connect and return if otherwise
	chatroom.Connect(sessionID, conn)
	if err != nil {
		conn.Close()
		fmt.Println("here2")
		return
	}

	// While connection is active
	for {
		// Open incoming message
		_, msg, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			fmt.Println("here3")
			return
		}

		// Parse the message's JSON
		var event MessageInboundEvent
		err = json.Unmarshal(msg, &event)
		if err != nil {
			conn.Close()
			fmt.Println("here5")
			return
		}

		// Forward message to Chatroom service
		err = chatroom.ReceiveMessage(sessionID, event)
		if err != nil {
			conn.Close()
			fmt.Println("here4")
			return
		}
	}
}

func main() {
	// Initialize channels
	chatroom.AddChannel("main")
	chatroom.AddChannel("off-topic")
	chatroom.AddChannel("new-members")

	// Router for static Svelte frontend
	svelte := http.FileServer(HTMLDir{http.Dir("./client/build/")})
	http.Handle("/", AuthMiddleware(http.StripPrefix("/", svelte)))

	// Websocket connection
	http.HandleFunc("/ws", Websocket)
	http.HandleFunc("/login.json", Register)

	http.ListenAndServe(":8080", nil)
}
