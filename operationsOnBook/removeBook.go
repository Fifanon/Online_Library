package bkop

import(
    "strconv"
	"online_library/modules/github.com/gorilla/mux"
	vars "online_library/varsAndFuncs"
	stct "online_library/structs"
	"net/http"
	dbconfig "online_library/config"
	searchbk "online_library/searchBook"
	s "online_library/session"

)

//RemoveBook **
func RemoveBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		vars.Tpl.ExecuteTemplate(w, "bookRemoving.html", nil)
		return
} 

//RemoveBookSearch **
func RemoveBookSearch(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		r.ParseForm()
		books := []stct.BookStruct{}
		var  found, errFound bool 

		if r.Form.Get("searchBy") == "isbn" {
			isbn, err := strconv.Atoi(r.Form.Get("value"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			books, found, errFound = searchbk.SearchByIsbn(isbn)
			if errFound {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	    if r.Form.Get("searchBy") == "title" {
			title:= r.Form.Get("value")
			books, found, errFound = searchbk.SearchByTitle(title)

		}
		if r.Form.Get("searchBy") == "author" {
			author := r.Form.Get("value")

			books, found, errFound = searchbk.SearchByTitle(author)
		}
		if !found {
				stct.Msg.BookExistsNot = "book does not exist"
				vars.Tpl.ExecuteTemplate(w, "bookRemoving.html", stct.Msg)
				stct.Msg.BookExistsNot = ""
				return
		}
		vars.Tpl.ExecuteTemplate(w, "bookToRemove.html", books)
		stct.Msg.BookExistsNot = ""
}

//RemoveBookprocessing **
func RemoveBookprocessing(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		isbn := params["isbn"]
		r.ParseForm()

		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		_, err = db.Query(`delete from book_instances where book_isbn = ?;`, isbn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Close()

		stct.Msg.Done = "Done"
		vars.Tpl.ExecuteTemplate(w, "bookRemoving.html", stct.Msg)

}
