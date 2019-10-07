package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

var users = map[string]string{
	os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
}

func Signin(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/index.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		username := r.FormValue("username")
		log.Printf("An auth attempt from %v!\n", username)
		password := r.FormValue("password")

		// Get the expected password from our in memory map
		expectedPassword, ok := users[username]

		// If a password exists for the given user
		// AND, if it is the same as the password we received, the we can move ahead
		// if NOT, then we return an "Unauthorized" status
		if !ok || bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password)) != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionManager.Put(r.Context(), "username", os.Getenv("USERNAME"))
		fmt.Fprintf(w, "Congratz, you authenticated!")

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")

	}

}

func Check(w http.ResponseWriter, r *http.Request) {

	if sessionManager.Exists(r.Context(), "username") {
		fmt.Fprintf(w, "Welcome!")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}

}

func Index(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "./static/index.html")
}

var sessionManager *scs.SessionManager

func main() {
	// Initialize a new session manager and configure the session lifetime, name and domain.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "SSSO_Cookie"
	sessionManager.Cookie.Domain = os.Getenv("DOMAIN")

	fmt.Printf("Will set auth for %v\n", os.Getenv("DOMAIN"))
	// "Signin" and "Welcome" are the handlers that we will implement
	mux := http.NewServeMux()
	mux.HandleFunc("/signin", Signin)
	mux.HandleFunc("/check", Check)
	mux.HandleFunc("/", Index)

	// start the server on port 8000
	fmt.Printf("Starting Simple Traefik SSO server\n")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", sessionManager.LoadAndSave(mux)))
}
