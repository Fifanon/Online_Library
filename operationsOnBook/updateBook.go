package bkop

import(
    "strconv"
	"github.com/gorilla/mux"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	searchbk "github.com/Fifanon/online_library/searchBook"
	s "github.com/Fifanon/online_library/session"
)

var	actualIsbn int

//UpdateBook **
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		vars.Tpl.ExecuteTemplate(w, "bookUpdating.html", nil)
		return
}

//UpdateBookSearch **
func UpdateBookSearch(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		books := []stct.BookStruct{}
		r.ParseForm()
		isbn, err := strconv.Atoi(r.Form.Get("value"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		books, found, errFound := searchbk.SearchByIsbn(isbn)
		if errFound {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !found {
			stct.Msg.BookExistsNot = "Book does not exist."
			vars.Tpl.ExecuteTemplate(w, "bookUpdating.html", stct.Msg)
			return
		}
		vars.Tpl.ExecuteTemplate(w, "bookUpdatingInput.html", books)
		return
}

//UpdateBookprocessing **
func UpdateBookprocessing(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	    params := mux.Vars(r)
	    var err error
	    actualIsbn, err = strconv.Atoi(params["isbn"])
        if err != nil{ 
           http.Error(w, err.Error(), http.StatusInternalServerError)
        }
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		qr, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area from book_instances where book_isbn = $1;`, actualIsbn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr.Next() {
			err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)		
		}
		r.ParseForm()
		isbn := r.Form.Get("isbn")
		title := r.Form.Get("title")
		author := r.Form.Get("authorname")
		pages:= r.Form.Get("pages")

		subjectArea := r.Form.Get("subject_area")

        if len(isbn) == 0 {
	       isbn = strconv.Itoa(stct.Bk.ISBN)
        } else{
		}
		if title == "" {
			title = stct.Bk.Title
		}
		if author == "" {
			author = stct.Bk.Author
		}
		if len(pages) == 0 {
			pages = strconv.Itoa(stct.Bk.Pages)
		}
		if subjectArea == "" {
			subjectArea = stct.Bk.Subject
		}
		newIsbn, err := strconv.Atoi(isbn)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		numOfpages, err := strconv.Atoi(pages)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = db.Query(`update book_instances set book_isbn = $1,book_title = $2,author_name = $3,pages = $4,subject_area = $5 where book_isbn = $6;`, newIsbn, title, author, numOfpages, subjectArea, actualIsbn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qrs, err := db.Query(`select book_isbn,book_title,author_name,pages,subject_area,number from book_instances where book_isbn = $1;`, newIsbn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for qrs.Next() {
			err = qrs.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if stct.Bk.Number > 0 {
				stct.Bk.Availability = "AVAILABLE"
			} else {
				stct.Bk.Availability = "NOT AVAILABLE"
			}
		}
		db.Close()
		vars.Tpl.ExecuteTemplate(w, "bookUpdated.html", stct.Bk)
		return
}

