package displaybk

import(
    "os"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	s "github.com/Fifanon/online_library/session"

)

//DisplayBiologyBooks **
func DisplayBiologyBooks(w http.ResponseWriter, r *http.Request) {

	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		vars.Subject = "Biology"

		books := []stct.BookStruct{}
		var enoughBr bool
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		borrqr, err := db.Query(`select count(*) from books_borrowed where member_email = $1;`, stct.User.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var fcount int
		var scount int

		for borrqr.Next() {
			err = borrqr.Scan(&fcount)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		if fcount < 5 {
			enoughBr = false
		} else {
			enoughBr = true
		}
		qr, err := db.Query(`select * from book_instances where subject_area = 'Biology';`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for qr.Next() {
			err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number, &stct.Bk.BookImageName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if stct.Bk.Number != 0 {
				stct.Bk.Availability = "AVAILABLE"
			} else {
				stct.Bk.Availability = "NOT AVAILABLE"
			}
			borrqr, err = db.Query(`select count(*) from books_borrowed where member_email = $1 and isbn = $2;`, stct.User.Email, stct.Bk.ISBN)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for borrqr.Next() {
				err = borrqr.Scan(&fcount)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			borrqr, err = db.Query(`select count(*) from tmp_borrow where mb_email = $1 and bk_isbn = $2;`, stct.User.Email, stct.Bk.ISBN)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for borrqr.Next() {
				err = borrqr.Scan(&scount)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			count := fcount + scount
			if enoughBr == true {
				stct.Bk.MsgToClient = "Not allowed to borrow!"
				stct.Bk.BorrAllowed = false
			} else {
				if count == 0 && stct.Bk.Number == 0 {
					stct.Bk.BorrAllowed = false
				} else if count == 0 && stct.Bk.Number != 0 {
					stct.Bk.BorrAllowed = true
				} else if count == 1 && stct.Bk.Number != 0 {
					stct.Bk.MsgToClient = "You have borrowed this book already !"
					stct.Bk.BorrAllowed = false
				} else if count == 1 && stct.Bk.Number == 0 {
					stct.Bk.MsgToClient = ""
					stct.Bk.BorrAllowed = false
				}
			}
			books = append(books, stct.Bk)
		}
		var lbemail string
		checkqr := db.QueryRow(`select l_email from librarian;`)
		err = checkqr.Scan(&lbemail)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if os.Getenv("EMAIL") == lbemail {
			vars.Tpl.ExecuteTemplate(w, "booksBiology.html", books)
			return
		}
		vars.Tpl.ExecuteTemplate(w, "m_booksBiology.html", books)
		return
}