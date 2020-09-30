package session

import (
    "os"
    "net/http"
    "github.com/gorilla/securecookie"
    stct "github.com/Fifanon/online_library/structs"
)
var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
	
//GetSession **
func GetSession(r *http.Request) (validated bool) {
    validated = true
    var name string = ""
    if cookie, err := r.Cookie("session-login"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("session-login", cookie.Value, &cookieValue); err == nil {
            name = cookieValue["session-login"]
        }
    }else{
        stct.Msg.LoginBefore = "Please fill the login form before you can assess"
        validated = false
        return validated
    }
    if name != os.Getenv("EMAIL") {
        stct.Msg.LoginBefore = "Please fill the login form before you can assess"
        validated = false
        return validated
    }
    return validated

}

//SetSession **
func SetSession(name string, response http.ResponseWriter) {
    value := map[string]string{
        "session-login": name,
    }
    if encoded, err := cookieHandler.Encode("session-login", value); err == nil {
        cookie := &http.Cookie{
            Name:  "session-login",
            Value: encoded,
            Path:  "/",
            MaxAge: 3600,
        }
        http.SetCookie(response, cookie)
    }else{
        panic(err)
    }

    return
}

//ClearSession **
func ClearSession(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "session-login",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}
 
 //ClearSessionHandler **
func ClearSessionHandler(response http.ResponseWriter, request *http.Request) {
    ClearSession(response)
    os.Unsetenv("EMAIL")
    http.Redirect(response, request, "/home", 302)
}