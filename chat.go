package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Users struct {
	Users []User `json:"users"`
}

type MessagesRequest struct {
	Username  string `json:"username"`
	OtherUser string `json:"otheruser"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	UsernameFrom string    `json:"usernameFrom"`
	UsernameTo   string    `json:"usernameTo"`
	Text         string    `json:"text"`
	Time         time.Time `json:"time"`
}

type User struct {
	ID          int       `json:"ID"`
	Username    string    `json:"username"`
	LastMessage time.Time `json:"time"`
	Online      bool      `json:"online"`
}

type Typing struct {
	UsernameFrom string `json:"usernameFrom"`
	UsernameTo   string `json:"usernameTo"`
	Typing       bool   `json:"typing"`
}

func chatTypingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data Typing
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "OK",
	})
	w.Write(jsonResponse)
	sendTyping(data)
}

func chatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data MessagesRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	ID1, err := getIDbyUsername(data.Username)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	ID2, err := getIDbyUsername(data.OtherUser)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	rows, err := statements["getMessages"].Query(ID1, ID2, ID2, ID1)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var messages MessagesResponse
	var IDpairs [][2]int
	for rows.Next() {
		var message Message
		var (
			IDFrom int
			IDTo   int
		)
		rows.Scan(&IDFrom, &IDTo, &message.Text, &message.Time)
		IDpairs = append(IDpairs, [2]int{IDFrom, IDTo})
		messages.Messages = append(messages.Messages, message)
	}
	rows.Close()
	for i := 0; i < len(messages.Messages); i++ {
		messages.Messages[i].UsernameFrom, err = getUsernamebyID(IDpairs[i][0])
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
			return
		}
		messages.Messages[i].UsernameTo, err = getUsernamebyID(IDpairs[i][1])
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
			return
		}
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(messages)
	w.Write(jsonResponse)
}

func chatUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data UsernameData
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	rows, err := statements["getAllUsers"].Query()
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var users Users
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Username)
		user.Online = isOnline(user.Username)
		users.Users = append(users.Users, user)
	}
	rows.Close()
	ID1, err := getIDbyUsername(data.Username)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	for i := 0; i < len(users.Users); i++ {
		ID2, err := getIDbyUsername(users.Users[i].Username)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
			return
		}
		rows, err := statements["getMessages"].Query(ID1, ID2, ID2, ID1)
		if err != nil {
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
			return
		}
		defer rows.Close()
		rows.Next()
		var (
			IDFrom int
			IDTo   int
			Text   string
		)
		rows.Scan(&IDFrom, &IDTo, &Text, &users.Users[i].LastMessage)
		rows.Close()
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(users)
	w.Write(jsonResponse)
}

func isOnline(username string) (found bool) {
	clients.Range(func(_, value interface{}) bool {
		if c, ok := value.(string); ok && c == username {
			found = true
		}
		return true
	})
	return
}

func chatNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "internal server error",
			})
			w.Write(jsonResponse)
		}
	}()
	var data Message
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	fromID, err := getIDbyUsername(data.UsernameFrom)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	toID, err := getIDbyUsername(data.UsernameTo)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Bad request",
		})
		w.Write(jsonResponse)
		return
	}
	data.Time = time.Now()
	_, err = statements["addMessage"].Exec(fromID, toID, data.Text, data.Time)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Message sent",
	})
	w.Write(jsonResponse)
	sendMessage(data)
}
