package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{}
	clients  = &sync.Map{}
)

type wsInput struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

type wsPost struct {
	Type string `json:"type"`
	Data Post   `json:"post"`
}

type wsComment struct {
	Type   string  `json:"type"`
	PostID int     `json:"postID"`
	Data   Comment `json:"comment"`
}

type wsStatus struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
}

type wsMessage struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

type wsTyping struct {
	Type         string `json:"type"`
	UsernameFrom string `json:"usernameFrom"`
	Typing       bool   `json:"typing"`
}

func wsConnection(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	reader(ws)
}

func reader(conn *websocket.Conn) {
	clients.Store(conn, "")
	defer clients.Delete(conn)
	defer conn.Close()
	for {
		_, incoming, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var data wsInput
		if err := json.Unmarshal([]byte(incoming), &data); err != nil {
			log.Println(err)
			return
		}
		switch data.Type {
		case "login":
			clients.Store(conn, data.Data["username"])
			sendStatus(data.Data["username"].(string), true)
			defer sendStatus(data.Data["username"].(string), false)
		case "logout":
			conn.Close()
			clients.Delete(conn)
			sendStatus(data.Data["username"].(string), false)
		}
	}
}

func sendPost(post Post) {
	data := wsPost{"post", post}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			c.WriteMessage(websocket.TextMessage, output)
		}
		return true
	})
}

func sendComment(postID int, comment Comment) {
	data := wsComment{"comment", postID, comment}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if c, ok := key.(*websocket.Conn); ok {
			c.WriteMessage(websocket.TextMessage, output)
		}
		return true
	})
}

func sendStatus(username string, online bool) {
	data := wsStatus{"status", username, online}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if value.(string) != "" {
			key.(*websocket.Conn).WriteMessage(websocket.TextMessage, output)
		}
		return true
	})
}

func sendMessage(message Message) {
	data := wsMessage{"message", message}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if value.(string) == message.UsernameFrom || value.(string) == message.UsernameTo {
			key.(*websocket.Conn).WriteMessage(websocket.TextMessage, output)
		}
		return true
	})
}

func sendTyping(typing Typing) {
	data := wsTyping{"typing", typing.UsernameFrom, typing.Typing}
	output, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	clients.Range(func(key, value interface{}) bool {
		if value.(string) == typing.UsernameTo {
			key.(*websocket.Conn).WriteMessage(websocket.TextMessage, output)
		}
		return true
	})
}
