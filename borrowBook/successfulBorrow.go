package borrowbk

import(
	"time"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	stct "github.com/Fifanon/online_library/structs"
	"net/http"
	dbconfig "github.com/Fifanon/online_library/config"
	s "github.com/Fifanon/online_library/session"
)


func processTime(mytime string) (deadline string, dayspassed float64) {

	now := time.Now()
	dayspassed = 0
	year := (int(mytime[0]) * 1000) + (int(mytime[1]) * 100) + (int(mytime[2]) * 10) + int(mytime[3])
	month := (int(mytime[5]) * 10) + int(mytime[6])
	day := (int(mytime[8]) * 10) + (int(mytime[9]) * 10)
	timePast := time.Date(year, time.Month(month), day, 0, 0, 0, 4815652, time.UTC)
	dur := timePast.AddDate(0, 0, 15)
	dura := now.Sub(dur)
	d, _ := time.ParseDuration(dura.String())
	h := d.Hours() / 24

	if h < 0 {
		diff := now.Sub(dur)
		dayspassed = diff.Hours() / 24
	} else if h == 0 {
		deadline = "today"
	} else {
		deadline = dur.Format("01-02-2006")
	}
	return deadline, dayspassed

}

//SuccBorrow **
func SuccBorrow(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		books := []stct.BookStruct{}
		var mytime string
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		qr, err := db.Query(`select book_isbn, book_title, author_name, pages,subject_area, number,b_imagename from
                       (book_instances join tmp_borrow on book_isbn = bk_isbn) where mb_email = $1;`, stct.User.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for qr.Next() {
			err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number, &stct.Bk.BookImageName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			stct.Bk.Fine = 0
			books = append(books, stct.Bk)
		}
		qr, err = db.Query(`select book_isbn,book_title, author_name, pages,subject_area,number,b_imagename,bowd_time,fine from
                       (book_instances join books_borrowed on book_isbn = bk_isbn) where member_email = $1;`, stct.User.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr.Next() {
			err = qr.Scan(&stct.Bk.ISBN, &stct.Bk.Title, &stct.Bk.Author, &stct.Bk.Pages, &stct.Bk.Subject, &stct.Bk.Number, &stct.Bk.BookImageName, &mytime, &stct.Bk.Fine)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			deadline, daysPassed := processTime(mytime)
			if daysPassed == 0 {
				stct.Bk.TimeLeft = deadline
			} else {
				stct.Bk.TimeLeft = "deadline is passed"
			}
		}
		vars.Tpl.ExecuteTemplate(w, "memberBooksborrowed.html", books)
		return
}