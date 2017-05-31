package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"

	"log"

	"sort"
)

var (
	master gameMaster
)

type player struct {
	Name           string `json:"name"`
	Score          int    `json:"score"`
	lastQuestionID int
}

type lobby struct {
	ID         int
	Course     course
	Module     module
	Players    []player
	questions  []question
	data       complete
	questionID int
	//Counter = antal frågor
	Counter int
	// TimeLimit = tidpunkten när nedräkningen är klar
	TimeLimit time.Time
	//TimeLeft = hur många sekunder kvar av TimeLimit
	TimeLeft int
	// lobbyState = 0 för lobby, 1 eller annat för game, och kanske 2 för end game ?
	lobbyState int
	// isTimeLimitSet = true om timeLimit har satts i när lobbyState är 0
	isTimeLimitSet bool
	die            int
}

type submittedAnswer struct {
	Answer  string
	Correct bool
}

type gameMaster struct {
	newID        int
	ActiveLobbys []lobby
	ActiveGames  []lobby
	dyingGames   []lobby
}

//------------------------------------------UNIVERSAL----------------------------------------------------------------

func removeLobby(lobbyActive *lobby, lobbys []lobby) []lobby {
	i := 0
	var remove int
	max := len(lobbys)
	for i < max {
		if lobbyActive.ID == lobbys[i].ID {
			remove = i
		}
		i = i + 1
	}
	lobbys = append(lobbys[:remove], lobbys[remove+1:]...)
	return lobbys
}

//från http://stackoverflow.com/questions/31487694/how-can-i-convert-string-to-integer-in-golang
func StrToInt(str string) (int, error) {
	nonFractionalPart := strings.Split(str, ".")
	return strconv.Atoi(nonFractionalPart[0])
}

func setTimeLeft(lobbyActive *lobby) {
	if isTimeSet(lobbyActive) {
		if isLobbysTimesUp(lobbyActive) {
			lobbyActive.TimeLeft = 0
			/*
				Without the else statement below "lobbyActive.TimeLeft = seconds" will allways be run,
				making the if statement above meaningless.
			*/
		} else {
			t0 := lobbyActive.TimeLimit
			t1 := time.Now()
			timeleft := t0.Sub(t1)
			seconds, err := StrToInt(timeleft.String())
			if err != nil {
				checkErr(err)
			}
			lobbyActive.TimeLeft = seconds
		}
	}
	return
}

func exit(res http.ResponseWriter, req *http.Request) {
	session := checkCookie(res, req)
	username := session.Values["username"].(string)
	out := "Removing " + username + " from lobby"
	log.Printf(out)
	activeLobby := lobbyFind(res, req)
	tmp := "lobby"
	if activeLobby == nil {
		activeLobby = gameFind(res, req)
		tmp = "game"
	}

	var remove int
	max := len(activeLobby.Players)
	/*
		if max == 1 {
			if tmp == "lobby" {
				removeLobby(activeLobby, master.ActiveLobbys)
			} else {
				removeLobby(activeLobby, master.ActiveGames)
			}
		}
	*/
	log.Println(tmp)
	i := 0
	for i < max {
		if username == activeLobby.Players[i].Name {
			remove = i
		}
		i = i + 1
	}
	activeLobby.Players = append(activeLobby.Players[:remove], activeLobby.Players[remove+1:]...)
}

func sendLobby(res http.ResponseWriter, lobbyActive *lobby) {
	output, err := json.Marshal(lobbyActive)
	if err != nil {
		fmt.Println("error: ", err)
	}
	log.Println(string(output))
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

//------------------------------------------------------------------------------------------------------------------------------
/*
.
.
.
*/
//-----------------------------------------------------LOBBY--------------------------------------------------------------------

func isLobbysTimesUp(lobbyActive *lobby) bool {
	if isTimeSet(lobbyActive) {
		t0 := lobbyActive.TimeLimit
		t1 := time.Now()
		timeDiff := t0.Sub(t1)
		return (timeDiff < 0) // returerar true om tiden har gått ut
	}
	return false
}

func isTimeSet(lobbyActive *lobby) bool {
	return lobbyActive.isTimeLimitSet
}

func isLobbyReady(lobbyActive *lobby) bool {
	if len(lobbyActive.Players) > 5 {
		setLobbyTimeLimit(lobbyActive, 0)
	}
	timeIsUp := isLobbysTimesUp(lobbyActive)
	enoughPlayers := checkEnoughPlayers(lobbyActive)
	if !enoughPlayers && timeIsUp {
		lobbyActive.isTimeLimitSet = false
		return false
	}
	return (enoughPlayers && timeIsUp) // true om tillräckligt med spelare och timern är slut
}

func splitReferer(referer string) []string {
	//VÄLDIGT HEMSK FUNKTION SOM SKA GÖRAS BÄTTRE OCH MER EFFEKTIV
	a := strings.Split(referer, "?")
	c := strings.Split(a[1], "&")
	d := strings.Split(c[0], "=")
	e := strings.Split(c[1], "=")
	var f []string
	f = append(f, d[1])
	f = append(f, e[1])
	return f
}

func lobbyFind(res http.ResponseWriter, req *http.Request) *lobby {
	var index *lobby
	i := 0
	session := checkCookie(res, req)
	if session == nil {
		return nil
	}
	tmp := splitReferer(req.Header.Get("Referer"))
	name := session.Values["username"].(string)
	course := tmp[0]
	courseID, _ := strconv.Atoi(course)
	moduleID, _ := strconv.Atoi(tmp[1])
	for i < len(master.ActiveLobbys) {
		if master.ActiveLobbys[i].Course.CourseID == courseID && master.ActiveLobbys[i].Module.ModuleID == moduleID {
			if master.ActiveLobbys[i].lobbyState == 1 && lobbyHasPlayer(&master.ActiveLobbys[i], name) {
				return &master.ActiveLobbys[i]
			} else if master.ActiveLobbys[i].lobbyState == 0 && len(master.ActiveLobbys[i].Players) < 700 { //700 player before lobby swaps to game
				index = &master.ActiveLobbys[i]
			}
		}
		i++
	}
	if index == nil {
		index = lobbyCreate(courseID, moduleID)
	}
	log.Println(name)
	if !lobbyHasPlayer(index, name) {
		log.Printf("Inserting a new player in lobby")
		lobbyInsertPlayer(index, name)
	}
	return index
}

func lobbyCreate(courseID int, moduleID int) *lobby {
	var lobbyActive lobby
	lobbyActive.Course = getCourseNameLocal(courseID)
	lobbyActive.Course.CourseID = courseID
	lobbyActive.Module = getModuleNameLocal(courseID, moduleID)
	lobbyActive.Module.ModuleID = moduleID
	lobbyActive.Counter = 2
	lobbyActive.lobbyState = 0
	master.newID = master.newID + 1
	lobbyActive.TimeLeft = 30
	lobbyActive.isTimeLimitSet = false
	lobbyActive.ID = master.newID
	master.ActiveLobbys = append(master.ActiveLobbys, lobbyActive)

	return &master.ActiveLobbys[(len(master.ActiveLobbys) - 1)]
}

func checkEnoughPlayers(lobbyActive *lobby) bool {
	amount := len(lobbyActive.Players)
	return (amount > 1)
}

func isLobbyFull(lobbyActive *lobby) bool {
	amount := len(lobbyActive.Players)
	return (amount >= 6) // false om amount < 6
}

func setLobbyToGame(lobbyActive *lobby) {
	lobbyActive.questions = getQuestions(lobbyActive.Course.CourseID, lobbyActive.Module.ModuleID)
	lobbyActive.Counter = len(lobbyActive.questions)
	lobbyActive.lobbyState = 1
	//setLobbyTimeLimit(lobbyActive, 30)
	master.ActiveGames = append(master.ActiveGames, *lobbyActive)
	//master.ActiveLobbys = removeLobby(lobbyActive, master.ActiveLobbys)
}

func getPlayerNames(lobbyActive *lobby) []string {
	var nameList []string
	fmt.Println(len(lobbyActive.Players))
	for i := 0; i < len(lobbyActive.Players); i++ {
		nameList[i] = lobbyActive.Players[i].Name
	}
	return nameList
}

func setLobbyTimeLimit(lobbyActive *lobby, limit time.Duration) {
	if !checkEnoughPlayers(lobbyActive) && lobbyActive.lobbyState == 0 {
		fmt.Println("Gick ej sätta timer")
		return // man ska inte kunna sätta timer då
	}
	lobbyActive.TimeLimit = time.Now().Add(limit * time.Second)
	setTimeLeft(lobbyActive)
	lobbyActive.isTimeLimitSet = true
}

func lobbyInsertPlayer(lobbyActive *lobby, username string) {
	var playerActive player
	playerActive.Name = username
	playerActive.Score = 0
	lobbyActive.Players = append(lobbyActive.Players, playerActive)
}

func lobbyPage(res http.ResponseWriter, req *http.Request) {
	log.Printf("lobbyPage")
	http.ServeFile(res, req, "./html/lobby.html")
}

func lobbyHasPlayer(activeLobby *lobby, name string) bool {
	for _, activePlayer := range activeLobby.Players {
		if activePlayer.Name == name {
			return true
		}
	}
	return false
}

func lobbyx(res http.ResponseWriter, req *http.Request) {
	log.Printf("lobbyx")
	lobbyActive := lobbyFind(res, req)
	if lobbyActive.lobbyState == 0 {
		log.Printf("lobbyState")
		if isLobbyReady(lobbyActive) {
			log.Printf("GameON!")
			setTimeLeft(lobbyActive)
			setLobbyToGame(lobbyActive)
			lobbyActive.die = 1
			//redirect sker i lobby.js
			return
		} else if !isTimeSet(lobbyActive) && checkEnoughPlayers(lobbyActive) {
			setLobbyTimeLimit(lobbyActive, 30)
		} else if isTimeSet(lobbyActive) && !checkEnoughPlayers(lobbyActive) {
			lobbyActive.TimeLeft = 30
			lobbyActive.isTimeLimitSet = false
		} else {
			setTimeLeft(lobbyActive)
		}
	} else {
		lobbyActive.die++
		log.Printf("GameON! not first")
		if lobbyActive.die == len(lobbyActive.Players) {
			master.ActiveLobbys = removeLobby(lobbyActive, master.ActiveLobbys)
		}
		return
	}
	sendLobby(res, lobbyActive)
}

//------------------------------------------------------------------------------------------------------------------------------
/*
.
.
.
*/
//-----------------------------------------GAME---------------------------------------------------------------------------------

func random(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func gameFind(res http.ResponseWriter, req *http.Request) *lobby {
	log.Printf("gameFind")
	i := 0
	session := checkCookie(res, req)
	if session == nil {
		return nil
	}
	name := session.Values["username"].(string)
	log.Printf(name)
	for i < len(master.ActiveGames) {
		for i, activeLobby := range master.ActiveGames {
			log.Printf("searching game")
			if activeLobby.lobbyState == 1 {
				log.Printf("lobbystate is 1")
				for _, activePlayer := range activeLobby.Players {
					if activePlayer.Name == name {
						return &master.ActiveGames[i]
					}
				}
			}
		}
		i++
	}
	return nil
}

func playerFind(res http.ResponseWriter, req *http.Request, lobbyActive *lobby) *player {
	session := checkCookie(res, req)
	name := session.Values["username"].(string)
	for i, playerActive := range lobbyActive.Players {
		if playerActive.Name == name {
			return &lobbyActive.Players[i]
		}
	}
	return nil
}

func gameServe(res http.ResponseWriter, req *http.Request) {
	log.Printf("GAMESERVE")
	lobbyActive := gameFind(res, req)
	if lobbyActive == nil {
		res.Write(nil)
		return
	}
	playerActive := playerFind(res, req, lobbyActive)
	log.Println(lobbyActive.Counter)
	var data complete
	if playerActive.lastQuestionID == lobbyActive.questionID && isLobbysTimesUp(lobbyActive) {
		if lobbyActive.Counter > 0 {
			setLobbyTimeLimit(lobbyActive, 15)
			lobbyActive.isTimeLimitSet = true
			max := len(lobbyActive.questions)
			remove := random(0, max)
			data.Question = lobbyActive.questions[remove]
			lobbyActive.questionID = data.Question.QuestionID
			playerActive.lastQuestionID = lobbyActive.questionID
			lobbyActive.Counter--
			lobbyActive.questions = append(lobbyActive.questions[:remove], lobbyActive.questions[remove+1:]...)
			data.Alts = getAlts(lobbyActive.questionID)
			lobbyActive.data = data
		} else {
			log.Printf("first else")
			res.Write(nil)
			return
		}
	} else {
		log.Printf("second else")
		data = lobbyActive.data
		playerActive.lastQuestionID = lobbyActive.questionID
	}
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func gamex(res http.ResponseWriter, req *http.Request) {
	log.Printf("gamex")
	lobbyActive := gameFind(res, req)
	setTimeLeft(lobbyActive)
	sendLobby(res, lobbyActive)
	log.Println(lobbyActive)
}

func sendCorrectAnswer(res http.ResponseWriter, req *http.Request) {
	lobbyActive := gameFind(res, req)
	qID := lobbyActive.questionID
	answer := getAnswer(qID)
	output, err := json.Marshal(answer)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func submitAnswer(res http.ResponseWriter, req *http.Request) {
	lobbyActive := gameFind(res, req)
	name := playerFind(res, req, lobbyActive)
	qID := lobbyActive.questionID
	var answer submittedAnswer
	answer.Answer = req.FormValue("answer")
	answer.Correct = checkAnswer(qID, answer.Answer)
	if answer.Correct {
		name.Score++
	}
	log.Println(answer.Correct)
	output, err := json.Marshal(answer.Correct)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)

}

//------------------------------------------------------------------------------------------------------------------------------
/*
.
.
.
*/
//---------------------------------------Score-------------------------------------------------------------------------------

type byScore []player

func (a byScore) Len() int           { return len(a) }
func (a byScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

func scoreLocal(res http.ResponseWriter, req *http.Request) {
	log.Printf("cleanup")
	cleanup()
	lobbyActive := gameFind(res, req)
	if lobbyActive != nil {
		lobbyActive.lobbyState = 2
		setLobbyTimeLimit(lobbyActive, 60)
	} else {
		lobbyActive = scoreFind(res, req)
	}
	sort.Sort(byScore(lobbyActive.Players))
	output, err := json.Marshal(lobbyActive.Players)
	if err != nil {
		fmt.Println("error: ", err)
	}
	log.Println(string(output))
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func cleanup() {
	for i, lobbyActive := range master.ActiveGames {
		if lobbyActive.lobbyState == 2 {
			setTimeLeft(&lobbyActive)
			log.Println(lobbyActive.isTimeLimitSet)
			log.Println(lobbyActive.TimeLeft)
			if isLobbysTimesUp(&lobbyActive) {
				master.ActiveGames = append(master.ActiveGames[:i], master.ActiveGames[i+1:]...)
				i--
			}
		}
	}
}

func scoreFind(res http.ResponseWriter, req *http.Request) *lobby {
	var index *lobby
	i := 0
	session := checkCookie(res, req)
	if session == nil {
		return nil
	}
	name := session.Values["username"].(string)
	log.Printf(name)
	for i < len(master.ActiveGames) {
		for i, activeLobby := range master.ActiveGames {
			if activeLobby.lobbyState == 2 {
				for _, activePlayer := range activeLobby.Players {
					if activePlayer.Name == name {
						return &master.ActiveGames[i]
					}
				}
			}
		}
	}
	return index
}

//------------------------------------------------------------------------------------------------------------------------------
/*
.
.
.
*/
//---------------------------------------Handlers-------------------------------------------------------------------------------

func testQuestions(res http.ResponseWriter, req *http.Request) {
	log.Printf("testQuestions")
	http.ServeFile(res, req, "./html/question.html")
}

func testScore(res http.ResponseWriter, req *http.Request) {
	log.Printf("testScore")
	http.ServeFile(res, req, "./html/score.html")
}

func coursePage(res http.ResponseWriter, req *http.Request) {
	log.Printf("coursePage")
	http.ServeFile(res, req, "./html/course.html")
}

func modulePage(res http.ResponseWriter, req *http.Request) {
	log.Printf("modulePage")
	http.ServeFile(res, req, "./html/module.html")
}

func errorPage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./html/error.html")
}

// logout handler
func logoutHandler(res http.ResponseWriter, req *http.Request) {
	clearSession(res, req)
	// The redirect doesn't do anything.
	http.Redirect(res, req, "/", 302)
}

func getUsername(res http.ResponseWriter, req *http.Request) {
	session := checkCookie(res, req)
	var username player
	username.Name = session.Values["username"].(string)
	output, err := json.Marshal(username)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func showData(res http.ResponseWriter, req *http.Request) {
	cleanup()
	output, err := json.Marshal(master)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func main() {
	//-------------------------------------Connect to database-----------------------------
	fmt.Println("Starting main")
	var err error
	db, err = sql.Open("mysql", "root:tieca@tcp(localhost:3306)/aceit")
	if err != nil {
		fmt.Println("main: error in handshake")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("main error in ping, you are not connected to the database")
		panic(err.Error())
	}

	//-----------------------------------------http------------------------------------------
	var unauthMux = http.NewServeMux()
	var authMux = http.NewServeMux()
	var router = http.NewServeMux()

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		activeMux := unauthMux
		session, err := store.Get(req, "logged")
		if err == nil {
			if _, ok := session.Values["authenticated"]; ok {
				activeMux = authMux
			}
		}
		handleFunc, _ := activeMux.Handler(req)
		handleFunc.ServeHTTP(res, req)
	})

	cssHandler := http.FileServer(http.Dir("css/"))
	imagesHandler := http.FileServer(http.Dir("images/"))
	jsHandler := http.FileServer(http.Dir("js/"))

	router.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	router.Handle("/images/", http.StripPrefix("/images/", imagesHandler))
	router.Handle("/js/", http.StripPrefix("/js/", jsHandler))

	unauthMux.HandleFunc("/signup", signupPage)
	unauthMux.HandleFunc("/", loginPage)
	unauthMux.HandleFunc("/error", errorPage)

	authMux.HandleFunc("/error", errorPage)
	authMux.HandleFunc("/", coursePage)
	authMux.HandleFunc("/course", coursePage)
	authMux.HandleFunc("/module", modulePage)
	authMux.HandleFunc("/logout", logoutHandler)

	authMux.HandleFunc("/getCourses", getCourses)
	authMux.HandleFunc("/getCourseName", getCourseName)
	authMux.HandleFunc("/getModules", getModules)
	authMux.HandleFunc("/getModuleName", getModuleName)
	authMux.HandleFunc("/getUsername", getUsername)

	authMux.HandleFunc("/exit", exit)

	authMux.HandleFunc("/data", showData)

	//---------------------Handlers for Lobby state---------------------------------

	authMux.HandleFunc("/lobby", lobbyPage)
	authMux.HandleFunc("/lobbyx", lobbyx)

	//---------------------Handlers for Game state---------------------------------

	authMux.HandleFunc("/question", testQuestions)
	authMux.HandleFunc("/gameServe", gameServe)
	authMux.HandleFunc("/game", testQuestions)
	authMux.HandleFunc("/gamex", gamex)
	authMux.HandleFunc("/getAnswer", sendCorrectAnswer)
	authMux.HandleFunc("/submitAnswer", submitAnswer)

	//----------------------Handlers for End state----------------------------------

	authMux.HandleFunc("/getLocalScore", scoreLocal)
	authMux.HandleFunc("/score", testScore)
	//http.HandleFunc("/getHighscore", gethighscore)

	//-----------------------------------------------------------------------------

	/*
		http.HandleFunc("/getPlayers", getplayers)
		http.HandleFunc("/getTimeLimit", gettimelimit)
		http.HandleFunc("/getTimeLeft", getlimeleft)
	*/

	http.ListenAndServe(":8080", context.ClearHandler(router))
}
