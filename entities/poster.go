package entities

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // Replace with your actual SQL driver
)

// /// Poster struct represents a poster entity in the database.///
type Poster struct {
	ID          string `json:"id"`
	Link_Poster string `json:"link_poster"`
	Note1       string `json:"note1"`
	Remark      string `json:"remark"`
}

// AllPoster retrieves all posters from the database.
func AllPoster() ([]Poster, error) {
	// PostgreSQL connection
	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("Select id, link_poster, note1, remark from poster;")
	if err != nil {

		log.Println(err)
		return nil, err

	}
	defer rows.Close()
	m := make([]Poster, 0)
	for rows.Next() {
		ml := Poster{}
		err := rows.Scan(&ml.ID, &ml.Link_Poster, &ml.Note1, &ml.Remark)
		if err != nil {
			return nil, err
		}
		m = append(m, ml)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return m, nil
}

// // UpdatePosterHandle updates an existing poster in the database.///
func UpdatePosterHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.FormValue("id")
	newNote1 := r.FormValue("note1")
	remark := r.FormValue("remark")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/invalid", http.StatusSeeOther)
		return
	}

	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get existing notes
	var existingNote1 string
	err = db.QueryRow("SELECT note1 FROM poster WHERE id = $1", id).Scan(&existingNote1)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/invalid", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch existing notes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cleanRemark := strings.TrimSpace(remark)

	// Get current timestamp in 12-hour format with AM/PM
	timestamp := time.Now().Format("01-02-2006 03:04:05 PM")

	// Check for "Approve Poster"
	if strings.EqualFold(cleanRemark, "Approve Poster") {
		existingNote1 = ""
	} else {
		// Append new note to note1 if present
		if newNote1 != "" {
			noteWithTime := "[" + timestamp + "] " + newNote1
			if existingNote1 != "" {
				existingNote1 += "\n\n" + noteWithTime
			} else {
				existingNote1 = noteWithTime
			}
		}

	}

	// Update database
	result, err := db.Exec(`
		UPDATE poster
		SET note1 = $1, remark = $2
		WHERE id = $3
	`, existingNote1, cleanRemark, id)

	if err != nil {
		http.Error(w, "Failed to update poster: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update result: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Redirect(w, r, "/invalid", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}

// // InsertPosterHandle handles the insertion of a new poster into the database.///
func InsertPosterHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	link_poster := r.FormValue("link_poster")

	// Default values for NOT NULL fields
	note1 := ""
	remark := ""

	// PostgreSQL connection
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
		INSERT INTO poster (id, link_poster, note1, remark)
		VALUES ($1, $2, $3, $4)
	`

	_, err = db.Exec(query, id, link_poster, note1, remark)
	if err != nil {
		http.Error(w, "Insert failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/allposter", http.StatusSeeOther)
}

// DeletePosterHandle handles the deletion of a poster from the database.///
func DeletePosterHandle(w http.ResponseWriter, r *http.Request) {
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
	query := `DELETE FROM poster WHERE id = $1`
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

	http.Redirect(w, r, "/allposter", http.StatusSeeOther)
}

// // UpdateLinkHandle updates the link of a poster in the database.///
func UpdateLinkHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/allposter", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	linkPoster := r.FormValue("link_poster")

	if id == "" || linkPoster == "" {
		http.Redirect(w, r, "/allposter", http.StatusSeeOther)
		return
	}

	//connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("DB error:", err)
		http.Redirect(w, r, "/allposter", http.StatusSeeOther)
		return
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE poster SET link_poster = $1 WHERE id = $2`, linkPoster, id)
	if err != nil {
		log.Println("Update error:", err)
	}

	http.Redirect(w, r, "/allposter", http.StatusSeeOther)
}
