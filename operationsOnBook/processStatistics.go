package bkop

import (
	"net/http"

	dbconfig "online_library/config"
	stct "online_library/structs"
	vars "online_library/varsAndFuncs"
	s "online_library/session"

)

//ProcessStatistics **
func ProcessStatistics(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		var actualStatistics stct.Statistics
		statisticsArr := []stct.Statistics{}
		var viewData stct.BorrowInfoarray
		var fstnumber, sndnumber int

		subjectAreas := []string{"Computer Engineering", "Mechanical Engineering", "Electronic Engineering", "Mathematics", "Physics", "Biology"}

		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}

		for _, subj := range subjectAreas {
			qr1, err := db.Query(`select number from book_instances where subject_area = ?;`, subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for qr1.Next() {
				err = qr1.Scan(&fstnumber)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			qr2, err := db.Query(`select count(*) from books_borrowed join book_instances on book_isbn = isbn and subject_area = ?;`, subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for qr2.Next() {
				err = qr2.Scan(&sndnumber)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			if subj == "Computer Engineering" {
				actualStatistics.NumOfCompEngBooks = fstnumber + sndnumber
			} else if subj == "Mechanical Engineering" {
				actualStatistics.NumOfMechEngBooks = fstnumber + sndnumber
			} else if subj == "Electronic Engineering" {
				actualStatistics.NumOfElectEngBooks = fstnumber + sndnumber
			} else if subj == "Mathematics" {
				actualStatistics.NumOfMathBooks = fstnumber + sndnumber
			} else if subj == "Physics" {
				actualStatistics.NumOfPhysicsBooks = fstnumber + sndnumber
			} else {
				actualStatistics.NumOfBiologyBooks = fstnumber + sndnumber
			}
			actualStatistics.NumOfBooks = actualStatistics.NumOfBooks + fstnumber + sndnumber
		}
		qr3, err := db.Query(`select count(*) from books_borrowed;`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr3.Next() {
			err = qr3.Scan(&actualStatistics.NumOfBorrowedBooks)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		qr4, err := db.Query(`select count(*) from books_borrowed where fine > 0;`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr4.Next() {
			err = qr4.Scan(&actualStatistics.NumOfMembOwingFine)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		qr5, err := db.Query(`select count(*) from members;`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr5.Next() {
			err = qr5.Scan(&actualStatistics.NumOfMembers)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		statisticsArr = append(statisticsArr, actualStatistics)
		viewData.Stat = statisticsArr
		vars.Tpl.ExecuteTemplate(w, "statistics.html", viewData)
}