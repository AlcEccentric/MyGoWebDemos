package main

import (
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func checkLog(foo func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if alreadyLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		foo(w, r)
	}

}

func isError(w http.ResponseWriter, r *http.Request, e error) bool {
	if e != nil {
		if e.Error() == incompleteError {
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return true
		}
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func checkNotLog(foo func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !alreadyLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		foo(w, r)
	}

}

func addCount(foo func(w http.ResponseWriter, r *http.Request, l *logInfo)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, u, logged := getUser(r)
		if logged && r.URL.Path != "/favicon.ico" {
			// fmt.Println(r.URL.Path)
			u.addVisitCount()
		}
		linfo := logInfo{
			User:   u,
			Logged: logged,
		}
		foo(w, r, &linfo)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) error {
	u := new(user)
	u.Email = r.FormValue("email")
	u.Fname = r.FormValue("firstname")
	u.Lname = r.FormValue("lastname")
	u.VisitCount = 0
	password := r.FormValue("password")

	if u.Email == "" || string(password) == "" {
		return infoImcomp{}
	}
	psw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = psw

	dbUsers[u.Email] = u
	// fmt.Println(dbUsers)
	return nil
}

func addSession(w http.ResponseWriter, r *http.Request) error {
	sid, err := uuid.NewV4()
	if err != nil {
		http.Error(w, "Cannot generate session ID", http.StatusInternalServerError)
		return err
	}

	c := mySetCookie(w, sid.String(), maxCookieLiveTime, "session")

	newSession := session{
		uid:       r.FormValue("email"),
		lastVisit: time.Now(),
	}
	dbSession[c.Value] = &newSession
	// fmt.Println(dbSession)
	return nil
}

func mySetCookie(w http.ResponseWriter, cookieValue string, cookieMaxAge int, cookieName string) *http.Cookie {
	c := http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		MaxAge:   cookieMaxAge,
		HttpOnly: true,
	}
	http.SetCookie(w, &c)
	return &c
}
