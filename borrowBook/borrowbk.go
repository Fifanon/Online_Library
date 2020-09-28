package borrowbk

import(
	"github.com/gorilla/mux"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	s "github.com/Fifanon/online_library/session"

)

//BorrowBook **
func BorrowBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		bookIsbn := mux.Vars(r)
		bIsbn := bookIsbn["isbn"]
		mEmail := stct.User.Email
		var brrNum, fine, totalfine int
		totalfine = 0

		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		qr, err := db.Query(`select count(*) from books_borrowed,tmp_borrow where mb_email = member_email and member_email = ?;`, stct.User.Email)
		if err != nil {
			panic(err)
		}
		for qr.Next() {
			err = qr.Scan(&brrNum)
			if err != nil {
				panic(err)
			}
		}
		if brrNum >= 5 {
			vars.Message = " You cannot borrow more than 5 books"
			vars.Tpl.ExecuteTemplate(w, "loggedIn.html", vars.Message)
			vars.Message = ""
			return
		}
		qr, err = db.Query(`select fine from books_borrowed where member_email = ?;`, stct.User.Email)
		if err != nil {
			panic(err)
		}
		for qr.Next() {
			err = qr.Scan(&fine)
			totalfine = totalfine + fine
			if err != nil {
				panic(err)
			}
		}
		if totalfine > 0 {
			vars.Message = " You cannot borrow because you owe fine"
			vars.Tpl.ExecuteTemplate(w, "loggedIn.html", vars.Message)
			vars.Message = ""
			return
		}
		temp, err := db.Prepare(`insert into tmp_borrow (bk_isbn,mb_email) values(?,?);`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = temp.Exec(&bIsbn, &mEmail)
		if err != nil {
			panic(err)
		}
		_, err = db.Query(`update book_instances set number=number-1 where book_isbn = ?;`, bIsbn)
		if err != nil {
			panic(err)
		}

		if vars.Subject == "Computer Engineering" {
			http.Redirect(w, r, "/sci-library/computer-eng-books", http.StatusSeeOther)
			return
		} else if vars.Subject == "Mechanical Engineering" {
			http.Redirect(w, r, "/sci-library/mechanical-eng-books", http.StatusSeeOther)
			return
		} else if vars.Subject == "Electronic Engineering" {
			http.Redirect(w, r, "/sci-library/electronic-eng-books", http.StatusSeeOther)
			return
		} else if vars.Subject == "Mathematics" {
			http.Redirect(w, r, "/sci-library/mathematics-books", http.StatusSeeOther)
			return
		} else if vars.Subject == "Physics" {
			http.Redirect(w, r, "/sci-library/physics-books", http.StatusSeeOther)
			return
		} else if vars.Subject == "Biology" {
			http.Redirect(w, r, "/sci-library/biology-books", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/sci-library/list-of-allbooks/ascending-order", http.StatusSeeOther)
			return
		}
}