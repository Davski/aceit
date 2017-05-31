package main

import (
	"net/http"

	"log"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	store = sessions.NewCookieStore(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
)

func checkCookie(res http.ResponseWriter, req *http.Request) *sessions.Session {
	session, err := store.Get(req, "logged")
	if err != nil {
		log.Printf("Problems with a cookie, make a new one")
		log.Println(err)
		clearSession(res, req)
		http.Redirect(res, req, "/index", 301)
		return nil
	}
	return session
}

func clearSession(res http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "logged")
	out := session.Values["username"].(string) + " logged out"
	log.Println(out)
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Options.MaxAge = -1
	session.Save(req, res)
}
