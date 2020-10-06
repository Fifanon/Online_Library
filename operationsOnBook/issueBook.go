package bkop

import(
	"time"
	"fmt"
	"github.com/gorilla/mux"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	s "github.com/Fifanon/online_library/session"
	gomail "github.com/Fifanon/online_library/gomail"

)


//IssuedBook **
func IssuedBook(w http.ResponseWriter, r *http.Request) {

	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
           panic(err)
		}
		qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number,b_imagename,m_imagename,m_firstname,m_lastname, m_email,m_status,m_address from tmp_borrow,members,book_instances
                       where book_isbn = bk_isbn and m_email = mb_email;`)
		if err != nil {
		}
		bookndMembs := []stct.BorrowInfo{}
		var bookndMem stct.BorrowInfo

		for qr.Next() {
			err = qr.Scan(&bookndMem.ISBN, &bookndMem.Title, &bookndMem.Author, &bookndMem.Pages, &bookndMem.Subject,
				&bookndMem.Number, &bookndMem.BookImageName, &bookndMem.ImageName, &bookndMem.FirstName, &bookndMem.LastName, &bookndMem.Email,&bookndMem.Status,&bookndMem.Address)
			if err != nil {
	           panic(err)
			}
			bookndMembs = append(bookndMembs, bookndMem)
		}
		db.Close() 

		if len(bookndMembs) == 0 {
			vars.Message = "NO BOOK REQUESTED"
			vars.Tpl.ExecuteTemplate(w, "NoData.html", vars.Message)
			vars.Message = ""
			return
		}
		vars.Tpl.ExecuteTemplate(w, "bookIssuing.html", bookndMembs)
		return
}

//SuccIssueBook **
func SuccIssueBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		email := params["email"]
		isbn := params["isbn"]
		fine := 0
		now := time.Now()
		deadline := now.AddDate(0, 0, 15)
		db, err := dbconfig.GetMySQLDb()
		if err != nil { 
           panic(err)
		}
		_, err = db.Exec(`insert into books_borrowed (isbn,member_email,fine,bowd_time,deadline)
              values($1,$2,$3,NOW(),$4);`,isbn,email,fine,deadline)
		if err != nil {
           panic(err)
		}
		qr, err := db.Query(`select author_name,book_title from book_instances where book_isbn = $1;`, isbn)
        if err != nil {
			panic(err)
		}
		var bookndMem stct.BorrowInfo
        for qr.Next() {
            err = qr.Scan(&bookndMem.Author,&bookndMem.Title)
            if err != nil {
				panic(err)            
			}
		}

		_, err = db.Query(`delete from tmp_borrow where mb_email = $1 and bk_isbn = $2;`, email, isbn)
		if err != nil {
           panic(err)
		}

		db.Close() 
		subject := "BOOK ISSUING"
		emailBody := fmt.Sprintf("You have been issued the book %s (%d) by %s.\n It is required that you return it in 2 weeks time.\n If deadline passed, you will have to pay $5 charge.\n", isbn,bookndMem.ISBN, bookndMem.Author)
		_,err = gomail.SendEmail(email,emailBody, subject)
		http.Redirect(w, r, "/sci-library/librarian/operations/issue-book", http.StatusSeeOther)
}



//DeleteBookRequest **
func DeleteBookRequest(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		email := params["email"]
		isbn := params["isbn"]

		r.ParseForm()

		db, err := dbconfig.GetMySQLDb()
		if err != nil {
           panic(err)
		}
		_, err = db.Query(`delete from tmp_borrow where mb_email = $1 and bk_isbn = $2;`, email, isbn)
		if err != nil {
           panic(err)
		}
		qr, err := db.Query(`select book_isbn,book_title,author_name from book_instances where book_isbn = $1;`, isbn)
        if err != nil {
			panic(err)
		}
		var bookndMem stct.BorrowInfo

        for qr.Next() {
            err = qr.Scan(&bookndMem.ISBN, &bookndMem.Title, &bookndMem.Author)
            if err != nil {
                     panic(err)            }
		}
		db.Close()

		subject := "BOOK BORROW REJECTED"
		emailBody := fmt.Sprintf("Your borrow of the book %s (%d) by %s, has been rejected. Contact the librarian on fifanonlesley@gmail to inquire about the reasons.\n", bookndMem.Title,bookndMem.ISBN, bookndMem.Author)
		_, err = gomail.SendEmail(email,emailBody, subject)
		http.Redirect(w, r, "/sci-library/librarian/operations/issue-book", http.StatusSeeOther)
}