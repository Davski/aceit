package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type course struct {
	CourseID int
	Course   string
}

type module struct {
	ModuleID int
	Module   string
	courseID int
}

type question struct {
	QuestionID int
	Question   string
	moduleID   int
	courseID   int
}

type alt struct {
	QuestionID int
	Alt1       string `json:"alt1"`
	Alt2       string `json:"alt2"`
	Alt3       string `json:"alt3"`
}

type complete struct {
	Question question
	Alts     alt
}

type ansswer struct {
	Answer string
}

func getCourses(res http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("select courseID, course from course")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var courses = []course{}
	for rows.Next() {
		var b = course{}
		if err = rows.Scan(&b.CourseID, &b.Course); err != nil {
			log.Println(err)
		}
		courses = append(courses, b)
	}
	output, err := json.Marshal(courses)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func getCourseName(res http.ResponseWriter, req *http.Request) {
	courseID := req.FormValue("courseID")
	row, err := db.Query("select course from course where courseID =?", courseID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var course course
	row.Next()
	if err = row.Scan(&course.Course); err != nil {
		log.Println(err)
		return
	}
	output, err := json.Marshal(course)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func getCourseNameLocal(courseID int) course {
	row, err := db.Query("select course from course where courseID =?", courseID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var course course
	row.Next()
	if err = row.Scan(&course.Course); err != nil {
		log.Println(err)
	}
	return course
}

func getModules(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Printf("getModules without posting courseID")
		http.Redirect(res, req, "/index", 300)
		return
	}

	courseID := req.FormValue("courseID")

	rows, err := db.Query("select moduleID, module from module where courseID =?", courseID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var modules = []module{}
	for rows.Next() {
		var b = module{}
		if err = rows.Scan(&b.ModuleID, &b.Module); err != nil {
			b.courseID, _ = strconv.Atoi(courseID)
			log.Println(err)
		}
		modules = append(modules, b)
	}
	output, err := json.Marshal(modules)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func getModuleName(res http.ResponseWriter, req *http.Request) {
	courseID := req.FormValue("courseID")
	moduleID := req.FormValue("moduleID")
	log.Printf("cID: %v, mID: %v", courseID, moduleID)
	row, err := db.Query("select module from module where courseID =? and moduleID =?", courseID, moduleID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var module module
	row.Next()
	if err = row.Scan(&module.Module); err != nil {
		log.Println(err)
		return
	}
	output, err := json.Marshal(module)
	if err != nil {
		fmt.Println("error: ", err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(output)
}

func getModuleNameLocal(courseID int, moduleID int) module {
	row, err := db.Query("select module from module where courseID =? and moduleID =?", courseID, moduleID)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var module module
	row.Next()
	if err = row.Scan(&module.Module); err != nil {
		log.Println(err)
	}
	return module
}

func getQuestions(courseID int, moduleID int) []question {

	rows, err := db.Query("select questionID, question from question where courseID =? and moduleID =?", courseID, moduleID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var questions = []question{}
	for rows.Next() {
		var b = question{}
		if err = rows.Scan(&b.QuestionID, &b.Question); err != nil {
			b.moduleID = moduleID
			b.courseID = courseID
			log.Println(err)
		}
		questions = append(questions, b)
	}
	return questions
}

func getAlts(questionID int) alt {

	rows, err := db.Query("select QuestionID, alt1, alt2, alt3 from alt where questionID = ?", questionID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var alts = []alt{}
	for rows.Next() {
		var b = alt{}
		rows.Scan(&b.QuestionID, &b.Alt1, &b.Alt2, &b.Alt3)
		alts = append(alts, b)
	}
	return alts[0]
}

func getAnswer(questionID int) string {
	rows, err := db.Query("select answer from answer where questionID = ?", questionID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var answerCorrect string
	rows.Scan(answerCorrect)
	return answerCorrect
}

func checkAnswer(questionID int, answer string) bool {
	log.Printf("questionID is: %d", questionID)
	rows, err := db.Query("select answer from answer where questionID = ?", questionID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var ansArray = []ansswer{}
	for rows.Next() {
		var b = ansswer{}
		if err := rows.Scan(&b.Answer); err != nil {
			log.Fatal(err)
		}
		ansArray = append(ansArray, b)
	}
	log.Printf("answerCorrect is: %s and your ans is %s: ", ansArray[0].Answer, answer)
	return ansArray[0].Answer == answer
}

/*
output, err := json.Marshal(questions)
if err != nil {
	fmt.Println("error: ", err)
}
res.Header().Set("Content-Type", "application/json")
res.Write(output)







type prepareS struct {
	course string
	module string
}

type preparedS struct {
	courseID   int
	moduleID   int
	questionID int
}

type activeQuestion struct {
	QuestionID int
	Question   string
}

type activeAlt struct {
	QuestionID int
	Alt1       string
	Alt2       string
	Alt3       string
}

type activeAnswer struct {
	QuestionID int
	Answer     string
}

type activeData struct {
	Questions []activeQuestion
	Alts      []activeAlt
	Answers   []activeAnswer
}

func api(db *sql.DB, course string, module string) []byte {

	var err error
	var question = []activeQuestion{}
	var alt = []activeAlt{}
	var answer = []activeAnswer{}
	var prepared preparedS
	var data activeData

	var inputFetch prepareS

	inputFetch.course = course
	inputFetch.module = module

	//----------------------------------------------------------------------------------------------

	err = db.QueryRow("select courseID from course where course = " + inputFetch.course).Scan(&prepared.courseID)
	switch {
		case err == sql.ErrNoRows:
		log.Printf("No course with name exists.")
		case err != nil:
		log.Fatal(err)
		default:
		//fmt.Println(prepared.courseID)
	}

	//----------------------------------------------------------------------------------------------

	err = db.QueryRow("select moduleID from module where module = "+inputFetch.module+"and courseID = ?", prepared.courseID).Scan(&prepared.moduleID)
	switch {
		case err == sql.ErrNoRows:
		log.Printf("No module with that name exists in the course " + inputFetch.course + ".")
		case err != nil:
		log.Fatal(err)
		default:
		//fmt.Println(prepared.courseID)
	}

	//----------------------------------------------------------------------------------------------

	rows, err := db.Query("select questionID, question from question where moduleID = ? and courseID = ?", prepared.moduleID, prepared.courseID)
	if err != nil {
		log.Fatal("DB query error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var b = activeQuestion{}
		if err := rows.Scan(&b.QuestionID, &b.Question); err != nil {
			log.Fatal(err)
		}
		question = append(question, b)
	}

	//----------------------------------------------------------------------------------------------

	rows, err = db.Query("select questionID, alt1, alt2, alt3 from alt where moduleID = ? and courseID = ?", prepared.moduleID, prepared.courseID)
	if err != nil {
		log.Fatal("DB query error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c = activeAlt{}
		if err := rows.Scan(&c.QuestionID, &c.Alt1, &c.Alt2, &c.Alt3); err != nil {
			log.Fatal(err)
		}
		alt = append(alt, c)
	}

	//----------------------------------------------------------------------------------------------

	rows, err = db.Query("select questionID, answer from answer where moduleID = ? and courseID = ?", prepared.moduleID, prepared.courseID)
	if err != nil {
		log.Fatal("DB query error: ", err)
	}
	defer rows.Close()
	for rows.Next() {
		var d = activeAnswer{}
		if err := rows.Scan(&d.QuestionID, &d.Answer); err != nil {
			log.Fatal(err)
		}
		answer = append(answer, d)
	}

	//----------------------------------------------gather all data-----------------------------

	data.Questions = question
	data.Alts = alt
	data.Answers = answer

	//-----------------------------------------------Json----------------------------------------

	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error: ", err)
	}
	//fmt.Println(string(output))
	return output
}
*/
