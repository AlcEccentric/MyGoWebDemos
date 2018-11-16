package main

import (
	"html/template"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type infoImcomp struct{}

func (i infoImcomp) Error() string {
	return incompleteError
}

type user struct {
	Email      string
	Fname      string
	Lname      string
	Password   []byte
	VisitCount int
}

func (u *user) addVisitCount() {
	u.VisitCount++
	return
}

type session struct {
	uid       string
	lastVisit time.Time
}

type logInfo struct {
	User   *user
	Logged bool
}

var incompleteError = "Info incomlete, please refresh the page."
var dbSession map[string]*session
var dbUsers map[string]*user
var maxCookieLiveTime = 60
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	dbSession = make(map[string]*session)
	dbUsers = make(map[string]*user)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if isError(w, r, r.ParseForm()) {
			return
		}
		// fmt.Println("Login email is:", r.FormValue("email"))
		// fmt.Println("Login pwd is:", r.FormValue("password"))
		u, ok := dbUsers[r.FormValue("email")]
		if !ok {
			http.Error(w, "Username do not match", http.StatusForbidden)
			return
		}

		pwd := r.FormValue("password")
		if err := bcrypt.CompareHashAndPassword(u.Password, []byte(pwd)); err != nil {
			http.Error(w, "Password do not match", http.StatusForbidden)
			return
		}

		if isError(w, r, addSession(w, r)) {
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
	isError(w, r, tpl.ExecuteTemplate(w, "login.html", nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		if isError(w, r, addUser(w, r)) || isError(w, r, addSession(w, r)) {

			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	isError(w, r, tpl.ExecuteTemplate(w, "signup.html", nil))
}

func logout(w http.ResponseWriter, r *http.Request) {

	c, _ := r.Cookie("session")
	// delete session
	delete(dbSession, c.Value)
	// delete cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func profile(w http.ResponseWriter, r *http.Request, l *logInfo) {
	// fmt.Println(*l)
	isError(w, r, tpl.ExecuteTemplate(w, "profile.html", l.User))
}

func root(w http.ResponseWriter, r *http.Request, l *logInfo) {
	isError(w, r, tpl.ExecuteTemplate(w, "index.html", l))
}

func main() {

	http.HandleFunc("/", addCount(root))
	http.HandleFunc("/signup", checkLog(signup))
	http.HandleFunc("/login", checkLog(login))
	http.HandleFunc("/profile", checkNotLog(addCount(profile)))
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)

}
