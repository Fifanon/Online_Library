package bkop

import (
	"net/http"
	 "time"
	 "strings"
	dbconfig "online_library/config"
	stct "online_library/structs"
	vars "online_library/varsAndFuncs"
	s "online_library/session"

)

var date []string
//ListOfBooksBorrowed **
func ListOfBooksBorrowed(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		var bookMem stct.BorrowInfo
		bookMems := []stct.BorrowInfo{}
		var lbemail string
		checkqr := db.QueryRow(`select l_email from librarian;`)
		err = checkqr.Scan(&lbemail)
		if err != nil {
			panic(err)
		}
		qr,err := db.Query("")
		if vars.Email == lbemail {
			qr, err = db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename,fine,deadline,m_firstname,m_lastname,m_email,m_imagename from books_borrowed, book_instances,members where isbn = book_isbn and m_email = member_email;`)
		}else{
			qr, err = db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename,fine,deadline,m_firstname,m_lastname,m_email,m_imagename from books_borrowed, book_instances,members where isbn = book_isbn and m_email = member_email and m_email = ?;`,vars.Email)
		}
		for qr.Next() {
			err = qr.Scan(&bookMem.ISBN, &bookMem.Title, &bookMem.Author, &bookMem.Pages, &bookMem.Subject, &bookMem.Number, &bookMem.BookImageName, &bookMem.Fine, &bookMem.Deadline, &bookMem.FirstName, &bookMem.LastName, &bookMem.Email,&bookMem.ImageName)
			if err != nil {
				panic(err)
			}
			bookMem.TimeLeft = subtractTime(bookMem.Deadline)
			bookMem.Deadline = date[0]
			if (bookMem.TimeLeft >= 0){
				bookMem.Deadline = "Passed deadline"
				bookMem.Fine = 5
				_, err = db.Query(`update books_borrowed set fine = ?`, bookMem.Fine)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}else{
				bookMem.TimeLeft = bookMem.TimeLeft *(-1)
			}
			bookMems = append(bookMems, bookMem)
		}
	    vars.Tpl.ExecuteTemplate(w, "memberBooksBorrowed.html", bookMems)

		return
}


func subtractTime(deadline string) int{
	date = strings.Fields(deadline)

	time1, err := time.Parse("2006-01-02", date[0])
	if err != nil {
		panic(err)
	}
	diff := time.Now().Sub(time1).Hours()
	daysLeft := int(diff/24)

    return daysLeft
}