package searchbook

import (
	"os"
	"net/http"
	"strconv"
	dbconfig "github.com/Fifanon/online_library/config"
	stct "github.com/Fifanon/online_library/structs"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	s "github.com/Fifanon/online_library/session"

)

//ProcessBookSearch **
func ProcessBookSearch(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	var bkfound, errEnc bool
	books := []stct.BookStruct{}
	db, err := dbconfig.GetMySQLDb()
	if err != nil {
		panic(err)
	}

	search := r.Form.Get("input")
	var lbemail string
	checkqr := db.QueryRow(`select l_email from librarian;`)
	err = checkqr.Scan(&lbemail)
	if err != nil {
		panic(err)
	}
	if r.Form.Get("searchBy") == "title" {
		books, bkfound, errEnc = SearchByTitle(search)
	} else if r.Form.Get("searchBy") == "isbn" {
		isbn, err := strconv.Atoi(search)
		if err != nil {
			stct.Msg.BookExistsNot = "Enter a valid Isbn"
			if os.Getenv("EMAIL") == lbemail {
				vars.Tpl.ExecuteTemplate(w, "operations.html",stct.Msg)
				stct.Msg.BookExistsNot = ""
				return
			}
			vars.Tpl.ExecuteTemplate(w, "loggedIn.html", stct.Msg)
			stct.Msg.BookExistsNot = ""
			return		}
		books, bkfound, errEnc = SearchByIsbn(isbn)
	} else {
		books, bkfound, errEnc = SearchByAuthor(search)
	}

	if errEnc == true {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !bkfound {
		stct.Msg.BookExistsNot = "There is no record of such book"
		if os.Getenv("EMAIL") == lbemail {
			http.Redirect(w, r, "/sci-library/librarian/operations", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/sci-library/welcome", http.StatusSeeOther)
		return
	} 
	if os.Getenv("EMAIL") == lbemail {
		vars.Tpl.ExecuteTemplate(w, "bookSearched.html", books)
		return
	}
	vars.Tpl.ExecuteTemplate(w, "m_bookSearched.html", books)
	return
}
