package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var (
	portHTTP string = "8080"
	fileDB   string = "./forum.db"
	db       *sql.DB
	reset    *bool
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Flags
	reset = flag.Bool("db-reset", false, "Reset database")
	flag.Parse()
	// DB
	dbInit()
	defer db.Close()

	// Static files, forum
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./web"))))
	log.Println("starting forum at http://localhost:" + portHTTP + "/")
	// Websocket
	http.HandleFunc("/ws", wsConnection)
	log.Println("starting websocket at ws://localhost:" + portHTTP + "/ws")
	// API
	http.HandleFunc("/api/user/register", userRegisterHandler)
	http.HandleFunc("/api/user/login", userLoginHandler)
	http.HandleFunc("/api/user/check", sessionCheckHandler)
	http.HandleFunc("/api/user/logout", userLogoutHandler)
	http.HandleFunc("/api/user/profile", userProfileHandler)
	http.HandleFunc("/api/post/submit", postNewHandler)
	http.HandleFunc("/api/post/get", postGetHandler)
	http.HandleFunc("/api/comment/submit", commentNewHandler)
	http.HandleFunc("/api/comment/get", commentGetHandler)
	http.HandleFunc("/api/chat/getusers", chatUsersHandler)
	http.HandleFunc("/api/chat/getmessages", chatMessagesHandler)
	http.HandleFunc("/api/chat/newmessage", chatNewHandler)
	http.HandleFunc("/api/chat/typing", chatTypingHandler)
	// Server
	log.Fatal(http.ListenAndServe(":"+portHTTP, nil))
}
