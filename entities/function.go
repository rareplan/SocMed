package entities

import (
	"database/sql"
	"log"
	"myproject/temp/config"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

func InitializeDatabase() *sql.DB {
	// Database connection details
	connStr := "host=localhost port=5432 user=postgres password=replan dbname=replan sslmode=disable"

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

type SessionData struct {
	Role string
}

var sessions = map[string]SessionData{}

// PARA MAG LOGIN //
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password") // Use `password` if plain-text; `password_hash` if already hashed

	connStr := "host=localhost port=5432 user=postgres password=replan dbname=replan sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var role string
	err = db.QueryRow(`SELECT role FROM users WHERE username=$1 AND password_hash=$2`, username, password).Scan(&role)
	if err != nil {
		http.Redirect(w, r, "/invalidlogin", http.StatusSeeOther)
		return
	}

	role = strings.ToLower(strings.TrimSpace(role)) // Always clean role

	// ✅ Save to session
	sessions[username] = SessionData{Role: role}

	// ✅ Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    username,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(60 * time.Minute),
	})

	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

// PARA HINDI BASTA BASTA PAMASOK KAHIT COPY ANG LINK//

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

// PARA SA DISPLAY DASHBOARD KASAMA ANG SESSION ROLE //
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

// PARA SA LOGOUT NG SESSION //
func logoutProcessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Remove cookie
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

// PARA SA VIEW YUNG CALENDAR //
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
