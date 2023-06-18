package main

import (
	"database/sql"
	"log"
	"os"
)

var statements = map[string]*sql.Stmt{}

func dbInit() {

	_, err := os.Stat(fileDB)

	if os.IsNotExist(err) {
		*reset = true
	}

	db, err = sql.Open("sqlite3", fileDB)
	if err != nil {
		log.Fatal(err)
	}
	if *reset {
		_, err := db.Exec(`
		DROP TABLE IF EXISTS user;
		DROP TABLE IF EXISTS session;
		DROP TABLE IF EXISTS post;
		DROP TABLE IF EXISTS category;
		DROP TABLE IF EXISTS post_category;
		DROP TABLE IF EXISTS comment;
		DROP TABLE IF EXISTS message;
		CREATE TABLE user (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			username VARCHAR UNIQUE NOT NULL,
			password VARCHAR NOT NULL,
			age INTEGER NOT NULL,
			gender VARCHAR NOT NULL,
			first_name VARCHAR NOT NULL,
			last_name VARCHAR NOT NULL,
			email VARCHAR UNIQUE NOT NULL
			);
		CREATE TABLE session (
			uuid VARCHAR PRIMARY KEY UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES user (id)
			);
		CREATE TABLE post (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES user (id),
			title VARCHAR NOT NULL,
			categories VARCHAR,
			text VARCHAR NOT NULL
			);
		CREATE TABLE comment (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			user_id INTEGER NOT NULL REFERENCES user (id),
			post_id INTEGER NOT NULL REFERENCES post (id),
			text VARCHAR NOT NULL
			);
		CREATE TABLE message (
			id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			from_id INTEGER NOT NULL REFERENCES user (id),
			to_id INTEGER NOT NULL REFERENCES user (id),
			text VARCHAR NOT NULL,
			time_sent DATETIME
			);`)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Println("DB reset")
	}
	for key, query := range map[string]string{
		"addUser":            `INSERT INTO user (username, age, gender, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?);`,
		"getAllUsers":        `SELECT id, username from user`,
		"getUserProfile":     `SELECT username, age, gender, first_name, last_name, email from user WHERE username=?`,
		"getUserbyID":        `SELECT username FROM user WHERE id = ?;`,
		"getUserID":          `SELECT id FROM user WHERE username = ?;`,
		"getUserCredentials": `SELECT username, password FROM user WHERE username = ? OR email = ?;`,
		"addSession":         `INSERT INTO session (uuid, user_id) VALUES (?, ?);`,
		"getSession":         `SELECT * FROM session WHERE uuid = ?;`,
		"getIDbyUUID":        `SELECT id FROM session INNER JOIN user ON user.id=user_id WHERE uuid = ?;`,
		"removeSession":      `DELETE FROM session WHERE uuid = ?;`,
		"addPost":            `INSERT INTO post (user_id, title, categories, text) VALUES (?, ?, ?, ?);`,
		"getPosts":           `SELECT post.id, username, title, categories, text FROM post INNER JOIN user ON user_id=user.id ORDER BY post.id DESC;`,
		"addComment":         `INSERT INTO comment (user_id, post_id, text) VALUES (?, ?, ?);`,
		"getComments":        `SELECT username, text FROM comment INNER JOIN user ON user_id = user.id WHERE post_id = ? ORDER BY comment.id DESC;`,
		"addMessage":         `INSERT INTO message (from_id, to_id, text, time_sent) VALUES (?, ?, ?, ?);`,
		"getMessages":        `SELECT from_id, to_id, text, time_sent FROM message WHERE (from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?) ORDER BY time_sent DESC;`,
	} {
		statements[key], err = db.Prepare(query)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
