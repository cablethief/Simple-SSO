package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]string{
	os.Getenv("USERNAME"): os.Getenv("PASSWORD"),
}

func Signin(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		http.ServeFile(res, req, "static/index.html")
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() err: %v", err)
			return
		}
		username := req.FormValue("username")
		log.Printf("An auth attempt from %v!\n", username)
		password := req.FormValue("password")

		// Get the expected password from our in memory map
		expectedPassword, ok := users[username]

		// If a password exists for the given user
		// AND, if it is the same as the password we received, the we can move ahead
		// if NOT, then we return an "Unauthorized" status
		if !ok || bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(password)) != nil {
			res.WriteHeader(http.StatusUnauthorized)
			log.Printf("Auth attempt from %v failed!\n", username)
			return
		}
		log.Printf("Auth attempt from %v succeeded!\n", username)
		sessionManager.Put(req.Context(), "username", os.Getenv("USERNAME"))
		fmt.Fprintf(res, "Congratz, you authenticated!")
		http.Redirect(res, req, req.Header.Get("CameFrom"), 302)

	default:
		fmt.Fprintf(res, "Sorry, only GET and POST methods are supported.")

	}

}

func Check(res http.ResponseWriter, req *http.Request) {
	if sessionManager.Exists(req.Context(), "username") {
		fmt.Fprintf(res, "Welcome!")
	} else {
		// res.Header().Set("ReferrerPolicy", "unsafe-url")
		// http.Redirect(res, req, "http://auth.local.host:8080/signin", http.StatusSeeOther)
		// res.WriteHeader(http.StatusUnauthorized)
	}
}

func Index(res http.ResponseWriter, req *http.Request) {
	if sessionManager.Exists(req.Context(), "username") {
		fmt.Fprintf(res, "Welcome!")
	} else {
		http.ServeFile(res, req, "./static/index.html")
	}
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
