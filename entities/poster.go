package entities

import (
	"database/sql"
	"log"

	"net/http"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3" // Replace with your actual SQL driver
)

type Poster struct {
	ID            string `json:"id"`
	Link_Poster   string `json:"link_poster"`
	Poster_number string `json:"poster_number"`
	Note1         string `json:"note1"`
	Note2         string `json:"note2"`
	Remark        string `json:"remark"`
}

// AllPoster retrieves all posters from the database//
func AllPoster() ([]Poster, error) {
	// PostgreSQL connection
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("Select id, link_poster, poster_number, note1, note2, remark from poster;")
	if err != nil {

		log.Println(err)
		return nil, err

	}
	defer rows.Close()
	m := make([]Poster, 0)
	for rows.Next() {
		ml := Poster{}
		err := rows.Scan(&ml.ID, &ml.Link_Poster, &ml.Poster_number, &ml.Note1, &ml.Note2, &ml.Remark)
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

// UpdatePosterHandle handles updating a poster's details in the database.
func UpdatePosterHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	posterNumber := r.FormValue("poster_number")
	note1 := r.FormValue("note1")
	note2 := r.FormValue("note2")
	remark := r.FormValue("remark")

	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get existing note1 and note2
	var existingNote1, existingNote2 string
	err = db.QueryRow("SELECT note1, note2 FROM poster WHERE id = $1", id).Scan(&existingNote1, &existingNote2)
	if err != nil {
		http.Error(w, "Failed to fetch existing notes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Clean remark input
	cleanRemark := remark
	if cleanRemark != "" {
		cleanRemark = strings.TrimSpace(cleanRemark)
	}

	// If remark is "Approve Poster" (case-insensitive), clear notes
	if strings.EqualFold(cleanRemark, "Approve Poster") {
		note1 = ""
		note2 = ""
	} else {
		// Retain existing notes if no new input
		if note1 == "" {
			note1 = existingNote1
		}
		if note2 == "" {
			note2 = existingNote2
		}
	}

	// Update record
	result, err := db.Exec(`
		UPDATE poster
		SET poster_number = $1, note1 = $2, note2 = $3, remark = $4
		WHERE id = $5
	`, posterNumber, note1, note2, cleanRemark, id)

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

// InsertPosterHandle handles the insertion of a new poster into the database.///
func InsertPosterHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	link_poster := r.FormValue("link_poster")

	// Default values for NOT NULL fields
	poster_number := ""
	note1 := ""
	note2 := ""
	remark := ""

	// PostgreSQL connection
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert query
	query := `
		INSERT INTO poster (id, link_poster, poster_number, note1, note2, remark)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = db.Exec(query, id, link_poster, poster_number, note1, note2, remark)
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
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
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

	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
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
