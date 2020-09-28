package member

import(
	"online_library/modules/github.com/gorilla/mux"
	s "online_library/session"
	vars "online_library/varsAndFuncs"
	stct "online_library/structs"
	"net/http"
	dbconfig "online_library/config"
)

//CancelMember **
func CancelMember(w http.ResponseWriter, r *http.Request) {

	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	vars.Tpl.ExecuteTemplate(w, "memberCancellation.html", nil)
	return
}



//CancelThisMember **
func CancelThisMember(w http.ResponseWriter, r *http.Request) {

	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	params := mux.Vars(r)

	db, err := dbconfig.GetMySQLDb()
	if err != nil {
		panic(err)
	}
	_, err = db.Query(`delete from  members where m_email = ?;`, params["email"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db.Close()
	vars.Tpl.ExecuteTemplate(w, "memberCancellation.html", stct.Msg)
	stct.Msg.Done = ""
}