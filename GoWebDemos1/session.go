package main

import (
	"net/http"
	"time"
)

var clearInterval = 30 * time.Second

func clearSessions() {

	for cookie, su := range dbSession {
		if time.Now().Sub(su.lastVisit) > clearInterval {
			delete(dbSession, cookie)
		}
	}

}

func getUser(r *http.Request) (*session, *user, bool) {
	c, err := r.Cookie("session")
	if err != nil {
		return nil, nil, false
	}
	su, ok := dbSession[c.Value]
	u := dbUsers[su.uid]
	return su, u, ok
}

func alreadyLoggedIn(r *http.Request) bool {
	su, _, ok := getUser(r)
	if ok {
		su.lastVisit = time.Now()
	}
	return ok
}
