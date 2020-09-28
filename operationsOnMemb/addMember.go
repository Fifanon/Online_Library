package member

import (
	"net/http"

	"online_library/modules/github.com/gorilla/mux"
	dbconfig "online_library/config"
	stct "online_library/structs"
	vars "online_library/varsAndFuncs"
	s "online_library/session"
)

var tmpMembers []stct.Users

//AddMember **
func AddMember(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		qr, err := db.Query(`select * from temporary_members;`)
		if err != nil {
			panic(err)
		}
		var tmpMember stct.Users
		tmpMembers := []stct.Users{}
		for qr.Next() {
			err = qr.Scan(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
			if err != nil {
				panic(err)
			}
			tmpMembers = append(tmpMembers, tmpMember)
		}
		if len(tmpMembers) == 0 {
			vars.Tpl.ExecuteTemplate(w, "NoData.html", "NO MEMBER TO ADD")
			return
		}
		vars.Tpl.ExecuteTemplate(w, "memberAdd.html", tmpMembers)
		return
}

//AddMembervalidate **
func AddMembervalidate(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		email := params["email"]
		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		qr, err := db.Query(`select * from temporary_members where mb_email = ?;`, email)
		if err != nil {
			panic(err)
		}
		var tmpMember stct.Users
		for qr.Next() {
			err = qr.Scan(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
			if err != nil {
				panic(err)
			}
		}
		temp, err := db.Prepare(`insert into members (m_firstname,m_lastname,m_email,m_address,m_telephone,m_password,m_status,m_imagename,m_signuptime)
              values(?,?,?,?,?,?,?,?,NOW());`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = temp.Exec(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
		if err != nil {
			panic(err)
		}
		qr, err = db.Query(`delete from temporary_members where mb_email = ?;`, email)
		if err != nil {
			panic(err)
		}
		qr, err = db.Query(`select * from temporary_members;`)
		if err != nil {
			panic(err)
		}
		for qr.Next() {
			err = qr.Scan(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
			if err != nil {
				panic(err)
			}
			tmpMembers = append(tmpMembers, tmpMember)
		}
		vars.Tpl.ExecuteTemplate(w, "memberAdd.html", tmpMembers)
		return
}

//DeleteRequest **
func DeleteRequest(w http.ResponseWriter, r *http.Request) {
	if	validated := s.GetSession(r);!validated{
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
		params := mux.Vars(r)
		email := params["email"]
		r.ParseForm()

		db, err := dbconfig.GetMySQLDb()
		if err != nil {
			panic(err)
		}
		_, err = db.Query(`delete from temporary_members where mb_email = ?;`, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Close()
		vars.Tpl.ExecuteTemplate(w, "memberAdd.html", tmpMembers)
}
