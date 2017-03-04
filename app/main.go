package main

import (
	"log"
	"net/http"
	"strconv"

	"google.golang.org/appengine"

	"golang.org/x/crypto/bcrypt"

	"cloud.google.com/go/datastore"
	"github.com/webapps/ridemerge"
	"golang.org/x/net/context"
)

var (
	ctx    context.Context
	client *datastore.Client
)

// server main
func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/register", register)
	http.HandleFunc("/post", createPost)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/dashboard/searchresults", searchResults)
	http.HandleFunc("/unsubscribe", unsubscribe)
	http.HandleFunc("/", index)
	log.Print("Listening on port 8080")
	appengine.Main()
}

// POST /register
// registers user.
func register(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	if session.LoggedIn {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	session.User = ridemerge.User{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Password:  string(hash[:]),
		Trips:     []ridemerge.Trip{},
	}
	var err error
	if _, err = ridemerge.UDB.GetUser(session.Email); err == nil {
		w.Write([]byte("false"))
		return
	}

	err = ridemerge.UDB.PutUser(session.User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cSession := &http.Cookie{
		Name:  "session-id",
		Value: session.Email,
	}
	http.SetCookie(w, cSession)
	w.Write([]byte("true"))
}

// POST /login
// login authenticates user.
func login(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	if session.LoggedIn {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	email := r.FormValue("email")
	var err error
	session.User, err = ridemerge.UDB.GetUser(email)
	if err != nil {
		w.Write([]byte("email"))
		return
	}
	password := r.FormValue("password")
	hash := []byte(session.Password)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		w.Write([]byte("password"))
		return
	}
	cSession := &http.Cookie{
		Name:  "session-id",
		Value: session.Email,
	}
	http.SetCookie(w, cSession)
	w.Write([]byte("true"))
}

// POST /logout
// logout deletes the session data and redirects user back to home.
func logout(w http.ResponseWriter, r *http.Request) {
	cSession, _ := r.Cookie("session-id")
	cSession = &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cSession)
	w.Write([]byte("success"))
}

// POST /logout
// logout deletes the session data and redirects user back to home.
func unsubscribe(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	var err error
	if session.LoggedIn {
		session.User, err = ridemerge.UDB.GetUser(session.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		password := r.FormValue("password")
		hash := []byte(session.Password)
		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		if err != nil {
			w.Write([]byte("false"))
			return
		}
		err = ridemerge.UDB.DeleteUser(session.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cSession, _ := r.Cookie("session-id")
		cSession = &http.Cookie{
			Name:   "session-id",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cSession)
		w.Write([]byte("true"))
		return
	}
	w.Write([]byte("false"))
}

func index(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	var err error
	if session.LoggedIn {
		session.User, err = ridemerge.UDB.GetUser(session.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	serveTemplate(w, session, "index", "navbar", "home", "login", "register", "post", "unsubscribe")
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	var err error
	if session.LoggedIn {
		session.User, err = ridemerge.UDB.GetUser(session.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		serveTemplate(w, session, "index", "navbar", "myposts", "login", "register", "post", "unsubscribe")
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func searchResults(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	var err error
	if !session.LoggedIn {
		http.Error(w, "Must be logged in", http.StatusUnauthorized)
		return
	} else {
		session.User, err = ridemerge.UDB.GetUser(session.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	searchText := r.URL.Query().Get("search")
	if searchText != "" {
		session.SearchResults, err = ridemerge.TDB.ListTripsWithDestination(searchText)
	} else {
		session.SearchResults, err = ridemerge.TDB.ListTrips()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serveTemplate(w, session, "index", "navbar", "searchbar", "searchposts", "login", "register", "post", "unsubscribe")
}

// POST /post
// createPost adds post to database.
func createPost(w http.ResponseWriter, r *http.Request) {
	session := ridemerge.GetSession(r)
	var err error
	if !session.LoggedIn {
		http.Error(w, "Must be logged in", http.StatusUnauthorized)
		return
	}
	event := r.FormValue("event")
	createdby := r.FormValue("createdby")
	origin := r.FormValue("origin")
	destination := r.FormValue("destination")
	dDate := r.FormValue("departureDate")
	dTime := r.FormValue("departureTime")
	var hasCar string
	if r.FormValue("hasCar") == "on" {
		hasCar = "yes"
	} else {
		hasCar = "no"
	}
	seats, _ := strconv.Atoi(r.FormValue("seatsAvailable"))

	trip := ridemerge.Trip{
		Event:       event,
		CreatedBy:   createdby,
		Origin:      origin,
		Destination: destination,
		DDate:       dDate,
		DTime:       dTime,
		HasCar:      hasCar,
		Seats:       seats,
	}
	if _, err = ridemerge.TDB.GetTrip(trip.Event); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	user, err := ridemerge.UDB.GetUser(session.Email)
	user.Trips = append(user.Trips, trip)
	err = ridemerge.UDB.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = ridemerge.TDB.PutTrip(trip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
