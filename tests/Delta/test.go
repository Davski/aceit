package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var course string
var module string

func startHTTPServer() *http.Server {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", homePage)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	// returning reference so caller can call Shutdown()
	return srv
}

func homePage(res http.ResponseWriter, req *http.Request) {

	http.ServeFile(res, req, "homePage.html")
	course = req.FormValue("course")
	module = req.FormValue("module")

	data := api(db, course, module)

	fmt.Printf(string(data))
	res.Write(data)
}

func main() {

	//-------------------------------------Connect to database-----------------------------
	fmt.Println("Starting main")
	db, err := sql.Open("mysql", "root:tieca@tcp(localhost:3306)/aceit")
	if err != nil {
		fmt.Println("main: error in handshake")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("main error in ping")
		panic(err.Error())
	}

	//-----------------------------------------http---------------------------------------

	srv := startHTTPServer()
	fmt.Println("main: listening for 30 seconds")
	time.Sleep(30 * time.Second)
	fmt.Printf("main: stopping server")
	if err2 := srv.Shutdown(nil); err2 != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	//------------------------------------------------------------------------------------

	fmt.Printf(course)
	fmt.Printf(module)

	//api(db, course, module)

	//------------------------------------------------------------------------------------

	log.Printf("main: done. exiting")

}
