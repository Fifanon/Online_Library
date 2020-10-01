package member

import (
	"net/http"
     "fmt"
	"github.com/gorilla/mux"
	dbconfig "github.com/Fifanon/online_library/config"
	gomail "github.com/Fifanon/online_library/gomail"
	stct "github.com/Fifanon/online_library/structs"
	vars "github.com/Fifanon/online_library/varsAndFuncs"
	s "github.com/Fifanon/online_library/session"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		qr, err := db.Query(`select * from temporary_members;`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var tmpMember stct.Users
		tmpMembers := []stct.Users{}
		for qr.Next() {
			err = qr.Scan(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		qr, err := db.Query(`select * from temporary_members where mb_email = $1;`, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var tmpMember stct.Users
		for qr.Next() {
			err = qr.Scan(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		temp, err := db.Prepare(`insert into members (m_firstname,m_lastname,m_email,m_address,m_telephone,m_password,m_status,m_imagename,m_signuptime)
              values($1,$2,$3,$4,$5,$6,$7,$8,NOW());`)
		if err != nil {
            panic(err)			
		}
		_, err = temp.Exec(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
		if err != nil {
			panic(err)			
		}
		qr, err = db.Query(`delete from temporary_members where mb_email = $1;`, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		qr, err = db.Query(`select * from temporary_members;`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		for qr.Next() {
			err = qr.Scan(&tmpMember.FirstName, &tmpMember.LastName, &tmpMember.Email, &tmpMember.Address, &tmpMember.PhoneNum, &tmpMember.Password, &tmpMember.Status, &tmpMember.ImageName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			tmpMembers = append(tmpMembers, tmpMember)
		}
		subject := "REGISTRATION CONFIRMATION"
		emailBody := fmt.Sprintf("Your registration at sci-library has been approved. click on the link below to login and get access to our great ressources.\n https://stormy-river-99671.herokuapp.com/home \n\n Regards.")
		_, err = gomail.SendEmail(stct.User.Email,emailBody, subject)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		_, err = db.Query(`delete from temporary_members where mb_email = $1;`, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.Close()
		subject := "REGISTRATION REJECTED"
		emailBody := fmt.Sprintf("Your registration at sci-library has been rejected. Contact the librarian at fifanonlesley@gmail to inquire about the reasons.\n")
		_, err = gomail.SendEmail(stct.User.Email,emailBody, subject)
		vars.Tpl.ExecuteTemplate(w, "memberAdd.html", tmpMembers)
		return
}
