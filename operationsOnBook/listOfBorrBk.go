package bkop

import (
	"fmt"
	"os"
	"net/http"
	 "time"
	 "strings"
	dbconfig "github.com/Fifanon/online_library/config"
	stct "github.com/Fifanon/online_library/structs"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	s "github.com/Fifanon/online_library/session"

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
		fmt.Println(os.Getenv("EMAIL"))
		if os.Getenv("EMAIL") == lbemail {
			qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename,fine,deadline,m_firstname,m_lastname,m_email,m_imagename,extract(days FROM (deadline - bowd_time)) from books_borrowed, book_instances,members where isbn = book_isbn and m_email = member_email;`)
			for qr.Next() {
				err = qr.Scan(&bookMem.ISBN, &bookMem.Title, &bookMem.Author, &bookMem.Pages, &bookMem.Subject, &bookMem.Number, &bookMem.BookImageName, &bookMem.Fine, &bookMem.Deadline, &bookMem.FirstName, &bookMem.LastName, &bookMem.Email,&bookMem.ImageName,&bookMem.TimeLeft)
				if err != nil {
					panic(err)
				}
				bookMem.Deadline = strings.Split(bookMem.Deadline, "T")[0]
				if bookMem.TimeLeft < 0{
					bookMem.Fine = 5
					bookMem.Deadline = bookMem.Deadline + "(Passed)"
					_, err = db.Query(`update books_borrowed set fine = $1;`,bookMem.Fine)
					if err != nil {
						panic(err)
					}
				}
				bookMems = append(bookMems, bookMem)
			}
		} else{
			qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename,fine,deadline,m_firstname,m_lastname,m_email,m_imagename, extract(days FROM (deadline - bowd_time)) from books_borrowed, book_instances,members where isbn = book_isbn and m_email = member_email and m_email = $1;`,os.Getenv("EMAIL"))
			for qr.Next() {
				err = qr.Scan(&bookMem.ISBN, &bookMem.Title, &bookMem.Author, &bookMem.Pages, &bookMem.Subject, &bookMem.Number, &bookMem.BookImageName, &bookMem.Fine, &bookMem.Deadline, &bookMem.FirstName, &bookMem.LastName, &bookMem.Email,&bookMem.ImageName,&bookMem.TimeLeft)
				if err != nil {
					panic(err)
				}
				bookMem.Deadline = strings.Split(bookMem.Deadline, "T")[0]
				if bookMem.TimeLeft < 0{
					bookMem.Fine = 5
					bookMem.Deadline = bookMem.Deadline + "(Passed)"
					_, err = db.Query(`update books_borrowed set fine = $1;`,bookMem.Fine)
					if err != nil {
						panic(err)
					}
				}				
				bookMems = append(bookMems, bookMem)
			}
		}
		if len(bookMems) == 0 {
			vars.Message = "No books borrowed"
			vars.Tpl.ExecuteTemplate(w, "NoData.html", vars.Message)
			vars.Message = ""
			return
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