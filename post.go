package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type PostRequest struct {
	UUID       string `json:"UUID"`
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Text       string `json:"text"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID         int    `json:"ID"`
	Username   string `json:"username"`
	Title      string `json:"title"`
	Categories string `json:"categories"`
	Text       string `json:"text"`
}

func postNewHandler(w http.ResponseWriter, r *http.Request) {
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
	var data PostRequest
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
	data.Categories = sanitizeCategories(data.Categories)
	ID, err := getIDbyUUID(data.UUID)
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = statements["addPost"].Exec(ID, data.Title, data.Categories, data.Text)
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
		"message": "Post created",
	})
	w.Write(jsonResponse)
	rows, err := statements["getPosts"].Query()
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var post Post
	rows.Next()
	rows.Scan(&post.ID, &post.Username, &post.Title, &post.Categories, &post.Text)
	rows.Close()
	sendPost(post)
}

func postGetHandler(w http.ResponseWriter, r *http.Request) {
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
	rows, err := statements["getPosts"].Query()
	if err != nil {
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	var posts Posts
	for rows.Next() {
		var post Post
		rows.Scan(&post.ID, &post.Username, &post.Title, &post.Categories, &post.Text)
		posts.Posts = append(posts.Posts, post)
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(posts)
	w.Write(jsonResponse)
}

func sanitizeCategories(data string) string {
	newdata := strings.Split(data, ",")
	returndata := ""
	for i, w := range newdata {
		w = strings.TrimSpace(w)
		if w != "" && i > 0 {
			returndata += (", " + w)
		} else {
			returndata += w
		}
		if returndata == "" {
			returndata = generateRandomEmojiSequence()
		}
	}
	return returndata
}

func generateRandomEmojiSequence() string {
	rounds := []string{"🔴", "🟠", "🟡", "🟢", "🔵", "🟣", "🟤", "⚫", "⚪"}
	// Shuffle the rounds using Fisher-Yates algorithm
	for i := len(rounds) - 1; i > 0; i-- {
		bi := big.NewInt(3)
		bj, err := rand.Int(rand.Reader, bi)
		if err != nil {
			log.Fatal(err)
		}
		// convert big.Int to int
		j := int(bj.Int64())
		rounds[i], rounds[j] = rounds[j], rounds[i]
	}

	// Join the shuffled rounds into a single string
	mixedRounds := strings.Join(rounds, " ")

	return mixedRounds
}
