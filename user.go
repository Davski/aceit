package main

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("Signup html")
		http.ServeFile(res, req, "html/signup.html")
		return
	}

	//-------------cookie-------------------
	session, _ := store.Get(req, "logged")
	//--------------------------------------

	username := req.FormValue("name")
	password := req.FormValue("password")
	passwordC := req.FormValue("Confirm password")

	if password != passwordC {
		log.Printf("the passwords didn't align")
		http.Redirect(res, req, "signup", 301)
		return
	}

	var user string

	err := db.QueryRow("SELECT name FROM user WHERE name=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account. 1", 500)
			return
		}

		_, err = db.Exec("INSERT INTO user(name, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account. 2", 500)
			return
		}

		//-------------cookie-------------------

		session.Values["authenticated"] = true
		session.Values["username"] = username
		session.Save(req, res)

		//--------------------------------------

		http.Redirect(res, req, "/course", 301)
		return
	case err != nil:
		http.Error(res, "Server error, unable to create your account. 3", 500)
		return
	default:
		http.Redirect(res, req, "/signup", 301)
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/index.html")
		return
	}

	//-------------cookie-------------------
	session, _ := store.Get(req, "logged")
	//--------------------------------------

	username := req.FormValue("name")
	password := req.FormValue("password")

	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT name, password FROM user WHERE name=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		log.Printf("No user with that name")
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		log.Printf("Password was wrong")
		http.Redirect(res, req, "/login", 301)
		return
	}

	//-------------cookie-------------------

	session.Values["authenticated"] = true
	session.Values["username"] = username
	session.Values["lobbyID"] = -1
	session.Save(req, res)

	//--------------------------------------
	http.Redirect(res, req, "/", 301)

}
