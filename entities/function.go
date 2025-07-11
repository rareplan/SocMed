package entities

import (
	"database/sql"
	"fmt"
	"log"
	"myproject/temp/config"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type SessionData struct {
	Role     string
	LoggedIn bool
}

var sessions = map[string]SessionData{}

// //////////////////////////////////////////// INITIALIZE DATABASE CONNECTION //////////////////////////////////////
func InitializeDatabase() *sql.DB {
	// Uncomment the following line to use a local database connection
	//connStr := "host=localhost port=5432 user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	// Check if connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	log.Println("Successfully connected to the database!")
	return db
}

// ///////////////////////////////////////////////////////////// LOGIN HANDLER //////////////////////////////////////////////////////
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Database connection
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get role from DB
	var role string
	err = db.QueryRow(`SELECT role FROM users WHERE username=$1 AND password_hash=$2`, username, password).Scan(&role)
	if err != nil {
		http.Redirect(w, r, "/invalidlogin", http.StatusSeeOther)
		return
	}

	// ✅ Check if user is already logged in
	if session, exists := sessions[username]; exists && session.LoggedIn {
		fmt.Println("Already logged in:", username) // for debugging
		http.Redirect(w, r, "/alreadylog", http.StatusSeeOther)
		return
	}

	// ✅ Clean role, save to session map
	role = strings.ToLower(strings.TrimSpace(role))
	sessions[username] = SessionData{
		Role:     role,
		LoggedIn: true,
	}

	// ✅ Set auth cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(60 * time.Minute),
	})

	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

// ///////////////////////////////////////////////////////// AUTHENTICATION MIDDLEWARE //////////////////////////////////////////////////////
func AuthMiddleware(allowedRoles []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		session, exists := sessions[cookie.Value]
		if !exists {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		for _, allowed := range allowedRoles {
			if session.Role == allowed {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Forbidden - Access Denied", http.StatusForbidden)
	}
}

// /////////////////////////////////////////////////////////// DASHBOARD HANDLER //////////////////////////////////////////////////////
func Dashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	session, ok := sessions[username]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Role flags for template
	data := struct {
		IsAdmin   bool
		IsUser    bool
		IsChecker bool
	}{
		IsAdmin:   session.Role == "admin",
		IsUser:    session.Role == "user",
		IsChecker: session.Role == "checker",
	}

	err = config.TPL.ExecuteTemplate(w, "/home", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ///////////////////////////////////////////////////////////// LOGOUT HANDLER //////////////////////////////////////////////////////
func logoutProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the auth_token cookie
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		username := cookie.Value

		// Unset session status
		if session, exists := sessions[username]; exists {
			session.LoggedIn = false
			sessions[username] = session
		}
	}

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ///////////////////////////////////////////////////////////// CALENDAR HANDLER //////////////////////////////////////////////////////
func Calendar(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	session, ok := sessions[username]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Role flags for template
	data := struct {
		IsAdmin   bool
		IsUser    bool
		IsChecker bool
	}{
		IsAdmin:   session.Role == "admin",
		IsUser:    session.Role == "user",
		IsChecker: session.Role == "checker",
	}

	err = config.TPL.ExecuteTemplate(w, "/calendar", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ///////////////////////////////////////////////////////////// NOTE HANDLER //////////////////////////////////////////////////////
func Note(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	session, ok := sessions[username]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Role flags for template
	data := struct {
		IsAdmin   bool
		IsUser    bool
		IsChecker bool
	}{
		IsAdmin:   session.Role == "admin",
		IsUser:    session.Role == "user",
		IsChecker: session.Role == "checker",
	}

	err = config.TPL.ExecuteTemplate(w, "/note", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ///////////////////////////////////////////////////////////// ACTIVITY HANDLER //////////////////////////////////////////////////////
func Act(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	session, ok := sessions[username]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Role flags for template
	data := struct {
		IsAdmin   bool
		IsUser    bool
		IsChecker bool
	}{
		IsAdmin:   session.Role == "admin",
		IsUser:    session.Role == "user",
		IsChecker: session.Role == "checker",
	}

	err = config.TPL.ExecuteTemplate(w, "/activity", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
