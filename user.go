package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	uuid "github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type signupData struct {
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type ProfileData struct {
	Username  string `json:"username"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type loginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UsernameData struct {
	Username string `json:"username"`
}

type UUIDData struct {
	UUID string `json:"UUID"`
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
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
	var profile ProfileData
	rows, err := statements["getUserProfile"].Query(data.Username)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "User not found",
		})
		w.Write(jsonResponse)
		return
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&profile.Username, &profile.Age, &profile.Gender, &profile.FirstName, &profile.LastName, &profile.Email)
	rows.Close()

	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(profile)
	w.Write(jsonResponse)
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	var data signupData
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

	onlyEnglishRegex := regexp.MustCompile(`^[a-zA-Z0-9]{2,15}$`)

	if !onlyEnglishRegex.MatchString(data.Username) {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid nickname: ` + data.Username + `
Nickname must only contain english letters and numbers.
Nickname must be between 2 and 15 characters long.`,
		})
		w.Write(jsonResponse)
		return
	}

	if data.Age < 1 || data.Age > 155 {
		w.WriteHeader(400)
		agestring := fmt.Sprintf("%d", data.Age)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid age: ` + agestring + `
Age must be between 1 and 155`,
		})
		w.Write(jsonResponse)
		return
	}

	if strings.ToLower(data.Gender) != "male" && strings.ToLower(data.Gender) != "female" {
		data.Gender = "undefined"
	}

	if len(data.FirstName) < 1 || len(data.FirstName) > 32 {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid first name length: ` + data.FirstName + `
First name must be between 1 and 32 characters long`,
		})
		w.Write(jsonResponse)
		return
	}

	if len(data.LastName) < 1 || len(data.LastName) > 32 {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid last name length: ` + data.LastName + `
Last name must be between 1 and 32 characters long`,
		})
		w.Write(jsonResponse)
		return
	}

	emailRegex := regexp.MustCompile(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,4})+$`)

	if !emailRegex.MatchString(strings.ToLower((data.Email))) {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid email: ` + data.Email + `
Email must be a valid email address`,
		})
		w.Write(jsonResponse)
		return
	}

	if len(data.Password) < 6 || len(data.Password) > 15 {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid password length: ` + data.Password + `
Password must be between 6 and 15 characters long`,
		})
		w.Write(jsonResponse)
		return
	}

	if !onlyEnglishRegex.MatchString(data.Password) {
		w.WriteHeader(400)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": `Invalid password: ` + data.Password + `
Password must only contain english characters and numbers`,
		})
		w.Write(jsonResponse)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	_, err = statements["addUser"].Exec(data.Username, data.Age, data.Gender, data.FirstName, data.LastName, data.Email, string(hash))
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: user.email" {
			log.Println(err.Error())
			w.WriteHeader(409)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "This email is already taken",
			})
			w.Write(jsonResponse)
			return
		}

		if err.Error() == "UNIQUE constraint failed: user.username" {
			log.Println(err.Error())
			w.WriteHeader(409)
			jsonResponse, _ := json.Marshal(map[string]string{
				"message": "This username is already taken",
			})
			w.Write(jsonResponse)
			return
		}

		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	UUID, err := createSession(data.Username)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"UUID":     UUID,
		"username": data.Username,
	})
	w.Write(jsonResponse)
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {

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
	var data loginData
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

	var username, hash string
	rows, err := statements["getUserCredentials"].Query(data.Username, data.Username)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid credentials",
		})
		w.Write(jsonResponse)
		return
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&username, &hash)
	rows.Close()

	if username == "" || hash == "" {
		w.WriteHeader(401)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid credentials",
		})
		w.Write(jsonResponse)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(data.Password))
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(401) // http status code 401 - http.StatusUnauthorized
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "Invalid credentials",
		})
		w.Write(jsonResponse)
		return
	}

	UUID, err := createSession(username)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"UUID":     UUID,
		"username": username,
	})
	w.Write(jsonResponse)
}

func sessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	var data UUIDData
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
	rows, err := statements["getSession"].Query(data.UUID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	defer rows.Close()
	if !rows.Next() {
		w.WriteHeader(200)
		jsonResponse, _ := json.Marshal(map[string]bool{
			"Exists": false,
		})
		w.Write(jsonResponse)
		return
	}
	rows.Close()
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]bool{
		"Exists": true,
	})
	w.Write(jsonResponse)
}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
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
	var data UUIDData
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
	_, err = statements["removeSession"].Exec(data.UUID)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(500)
		jsonResponse, _ := json.Marshal(map[string]string{
			"message": "internal server error",
		})
		w.Write(jsonResponse)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")
	w.WriteHeader(200)
	jsonResponse, _ := json.Marshal(map[string]string{
		"message": "Session deleted",
	})
	w.Write(jsonResponse)
}

func createSession(username string) (UUID string, err error) {
	random, _ := uuid.NewV4()
	UUID = random.String()
	ID, err := getIDbyUsername(username)
	if err != nil {
		return "", err
	}
	_, err = statements["addSession"].Exec(UUID, ID)
	if err != nil {
		return "", err
	}
	return UUID, nil
}

func getIDbyUsername(username string) (ID int, err error) {
	rows, err := statements["getUserID"].Query(username)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&ID)
	rows.Close()
	return ID, nil
}

func getIDbyUUID(UUID string) (ID int, err error) {
	rows, err := statements["getIDbyUUID"].Query(UUID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&ID)
	rows.Close()
	return ID, nil
}

func getUsernamebyID(ID int) (username string, err error) {
	rows, err := statements["getUserbyID"].Query(ID)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	rows.Next()
	rows.Scan(&username)
	rows.Close()
	return username, nil
}
