package main

import (
	"database/sql"
	"fmt"
	"log"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

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
