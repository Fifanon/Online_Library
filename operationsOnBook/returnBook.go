package bkop

import(
	"online_library/modules/github.com/gorilla/mux"
	vars "online_library/varsAndFuncs"
	stct "online_library/structs"
	"net/http"
	dbconfig "online_library/config"
	s "online_library/session"

)

//ReturnBook **
func ReturnBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		vars.Tpl.ExecuteTemplate(w, "bookReturning.html", nil)
		return
}

//ReturnBookMemberSearch **
func ReturnBookMemberSearch(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		r.ParseForm()
		var email string
		booknBorrs := []stct.BorrowInfo{}

		if r.Form.Get("email") == "" {

		} else {
			email = r.Form.Get("email")
		}
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}

		qr, err := db.Query(`select bb.isbn,m.m_imagename,m.m_firstname,m.m_lastname, m.m_email,bb.fine from books_borrowed bb join members as m on m.m_email = bb.member_email where bb.member_email = ?;`, email)
		if err != nil {
			panic(err)
		}
		for qr.Next() {
			err = qr.Scan(&stct.BooknBorr.ISBN, &stct.BooknBorr.ImageName, &stct.BooknBorr.FirstName, &stct.BooknBorr.LastName, &stct.BooknBorr.Email, &stct.BooknBorr.Fine)
			if err != nil {
				panic(err)
			}

			isbn := stct.BooknBorr.ISBN
			qrb, err := db.Query(`select book_title,author_name,b_imagename, subject_area from book_instances where book_isbn = ?;`, isbn)
			if err != nil {
				panic(err)
			}

			for qrb.Next() {
				err = qrb.Scan(&stct.BooknBorr.Title, &stct.BooknBorr.Author, &stct.BooknBorr.BookImageName,&stct.BooknBorr.Subject)
					if err != nil {
						panic(err)
					}
			}
			booknBorrs = append(booknBorrs, stct.BooknBorr)
			
		}
		if len(booknBorrs) == 0{
			stct.Msg.BookExistsNot = "member has not borrowed any book"
			vars.Tpl.ExecuteTemplate(w, "bookReturning.html", stct.Msg)
		}
		db.Close()
		vars.Tpl.ExecuteTemplate(w, "booksToReturn.html", booknBorrs)
		return
}

//SuccReturnBook **
func SuccReturnBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		isbn := params["isbn"]
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		_, err = db.Query(`delete from books_borrowed where isbn = ?;`, isbn)
		if err != nil {
			panic(err)
		}
		_, err = db.Query(`update book_instances set number = number+1 where book_isbn = ?;`, isbn)
		if err != nil {
			panic(err)
		}
		stct.Msg.Done = "Done"
		vars.Tpl.ExecuteTemplate(w, "bookReturning.html", stct.Msg)
		return
}