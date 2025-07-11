package entities

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // Replace with your actual SQL driver
)

// // User struct represents a user entity in the database.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password_hash"`
	Role     string `json:"role"`
}

// AllUsers retrieves all users from the database.//
func AllUsers() ([]User, error) {
	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, username, password_hash, role FROM users;")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// DeleteUserHandle handles the deletion of a user from the database.//
func DeleteUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")

	// PostgreSQL connection
	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete query//
	query := `DELETE FROM users WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Delete failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete result: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		return
	}

	http.Redirect(w, r, "/useraccess", http.StatusSeeOther)
}

// InsertUserHandle handles the insertion of a new user into the database.//
func InsertUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert query
	query := `
		INSERT INTO users (id, username, password_hash, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(query, id, username, password, role)
	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/useraccess", http.StatusSeeOther)
}

// UpdateUserHandle handles the update of an existing user in the database.//
func UpdateUserHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	newUsername := r.FormValue("username")
	newPassword := r.FormValue("password")
	newRole := r.FormValue("role")

	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get current values
	var currentUsername, currentPassword, currentRole string
	err = db.QueryRow(`SELECT username, password_hash, role FROM users WHERE id = $1`, id).Scan(&currentUsername, &currentPassword, &currentRole)
	if err != nil {
		http.Error(w, "User not found: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use old values if new ones are empty
	if newUsername == "" {
		newUsername = currentUsername
	}
	if newPassword == "" {
		newPassword = currentPassword
	}
	if newRole == "" {
		newRole = currentRole
	}

	// Update user
	_, err = db.Exec(`
		UPDATE users
		SET username = $2, password_hash = $3, role = $4
		WHERE id = $1
	`, id, newUsername, newPassword, newRole)

	if err != nil {
		http.Error(w, "Update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/useraccess", http.StatusSeeOther)
}

// InsertData handles the insertion of a new user into the database via a form submission.//
func InsertData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	username := r.FormValue("username")
	password := r.FormValue("password")
	role := r.FormValue("role")

	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert query
	query := `
		INSERT INTO users (id, username, password_hash, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(query, id, username, password, role)
	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
