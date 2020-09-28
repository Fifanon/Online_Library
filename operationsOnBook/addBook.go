package bkop

import(
	"os"
	"database/sql"
	"io/ioutil"
	"fmt"
	"strconv"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	s "github.com/Fifanon/online_library/session"

)

//AddBook **
func AddBook(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		vars.Tpl.ExecuteTemplate(w, "bookAdding.html", nil)
		return
}

//BookAdding **
func BookAdding(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			stct.Msg.Any = "Upload a book cover please"
			vars.Tpl.ExecuteTemplate(w, "bookAdding.html", stct.Msg)
			stct.Msg.Any = ""
		}
	   vars.CoverFilename = handler.Filename
		vars.Coverfilebytes, err = ioutil.ReadAll(file)
		if err != nil {
			stct.Msg.Any = "Choose a file please"
			vars.Tpl.ExecuteTemplate(w, "bookAdding.html", stct.Msg)
			panic(err)
		}
		defer file.Close()
		subjectAreas := [6]string{"Mathematics", "Physics", "Biology", "Computer Engineering", "Electronic Engineering", "Mechanical Engineering"}
		
		r.ParseForm()
		isbn := r.FormValue("isbn")
		if err !=nil {
			panic(err)
		}
		stct.Bk.ISBN,err = strconv.Atoi(isbn)

		stct.Bk.Title = r.FormValue("title")
		stct.Bk.Author = r.FormValue("authorname")
		pages:= r.FormValue("pages")
		stct.Bk.Pages,err = strconv.Atoi(pages)
		if err !=nil {
			panic(err)
		}
		number := r.FormValue("number")
		if err !=nil {
			panic(err)
		}
		stct.Bk.Number,err = strconv.Atoi(number)
		if err !=nil {
			panic(err)
		}

		stct.Bk.Subject = r.FormValue("subject_area")
		stct.Bk.BookImageName = vars.CoverFilename

	    fmt.Println(stct.Bk.ISBN)

		db, err := dbconfig.GetMySQLDb()
			if err != nil {
				panic(err)
			}
			qr := db.QueryRow(`select book_title from book_instances where book_isbn = ?;`, isbn)
			err = qr.Scan(&stct.Bk.Title)
			if err == nil {
				stct.Msg.Any = "This book already exists"
				http.Redirect(w, r, "/sci-library/librarian/operations/add-book", http.StatusSeeOther)
				return
			}

			if err != nil {
				if err == sql.ErrNoRows {

					tmp, err := db.Prepare(`INSERT into book_instances (book_isbn,book_title,author_name,pages,subject_area,number,b_imagename)
                                     values(?,?,?,?,?,?,?)`)
					if err != nil {
						panic(err)
					}
					_, err = tmp.Exec(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject,&stct.Bk.Number, &stct.Bk.BookImageName)
					if err != nil {
						panic(err)
					}
					for i := 0; i < 6; i++ {
						if subjectAreas[i] == stct.Bk.Subject{
							f, err := os.OpenFile("./project_files//public/subj-img"+subjectAreas[i]+"/"+vars.CoverFilename, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								panic(err)
							}
							f, err = os.OpenFile("./project_files//public/subj-img/books/"+vars.CoverFilename, os.O_WRONLY|os.O_CREATE, 0666)
							if err != nil {
								panic(err)
							}
							f.Write(vars.Coverfilebytes)
						}
					}
					vars.Tpl.ExecuteTemplate(w, "bookAddedSucc.html", stct.Bk)
					vars.GotFile = false
					return
				}
				panic(err)
			}
}
