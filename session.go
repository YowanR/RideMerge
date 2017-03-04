package ridemerge

import "net/http"

// getSession returns the current session if it exists, otherwise an empty
// session is returned.
func GetSession(r *http.Request) Session {
	cUser, err := r.Cookie("session-id")
	var session Session
	if err == nil {
		session.Email = cUser.Value
		session.LoggedIn = true
	}
	return session
}
