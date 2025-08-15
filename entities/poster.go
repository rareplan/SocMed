package entities

import (
	"database/sql"
	"io"
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
	Image_data  []byte `json:"image_data"` // Optional field for image data
	Remark      string `json:"remark"`
}

// AllPoster retrieves all posters from the database.
func AllPoster() ([]Poster, error) {
	// PostgreSQL connection
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	// Live Connection
	//connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("Select id, link_poster, note1, image_data, remark from poster;")
	if err != nil {

		log.Println(err)
		return nil, err

	}
	defer rows.Close()
	m := make([]Poster, 0)
	for rows.Next() {
		ml := Poster{}
		err := rows.Scan(&ml.ID, &ml.Link_Poster, &ml.Note1, &ml.Image_data, &ml.Remark)
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

	// ====== Read uploaded image ======
	file, _, err := r.FormFile("image_data")
	var imageBytes []byte
	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "Failed to read image: "+err.Error(), http.StatusBadRequest)
		return
	}
	if file != nil {
		defer file.Close()
		imageBytes, err = io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read image content: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/invalid", http.StatusSeeOther)
		return
	}

	connStr := "host=localhost port=5432 user=postgres password=replan dbname=replan sslmode=disable"
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
	timestamp := time.Now().Format("01-02-2006 03:04:05 PM")

	// Update note logic
	if strings.EqualFold(cleanRemark, "Approve Poster") {
		existingNote1 = ""
	} else {
		if newNote1 != "" {
			noteWithTime := "[" + timestamp + "] " + newNote1
			if existingNote1 != "" {
				existingNote1 += "\n\n" + noteWithTime
			} else {
				existingNote1 = noteWithTime
			}
		}
	}

	// ====== Update logic ======
	if len(imageBytes) > 0 {
		// may bagong image → update kasama image_data
		_, err = db.Exec(`
            UPDATE poster
            SET note1 = $1, remark = $2, image_data = $3
            WHERE id = $4
        `, existingNote1, cleanRemark, imageBytes, id)
	} else {
		// walang bagong image → update note1 at remark lang
		_, err = db.Exec(`
            UPDATE poster
            SET note1 = $1, remark = $2
            WHERE id = $3
        `, existingNote1, cleanRemark, id)
	}

	if err != nil {
		http.Error(w, "Failed to update poster: "+err.Error(), http.StatusInternalServerError)
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
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	// Live Connection
	//connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
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
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	// Live Connection
	//connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
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
	// PostgreSQL connection
	connStr := "user=postgres password=replan dbname=replan sslmode=disable"
	// Live Connection
	//connStr := "host=dpg-d1n2fkuuk2gs739eu39g-a.oregon-postgres.render.com port=5432 user=replan_sz89_user password=xkMmzaTtoqm9NouEyVaXWMZGgsdamovb dbname=replan_sz89 sslmode=require"
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

func ServeImage(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/image/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=replan dbname=replan sslmode=disable")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var imageData []byte
	err = db.QueryRow("SELECT image_data FROM poster WHERE id = $1", id).Scan(&imageData)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	contentType := http.DetectContentType(imageData)
	w.Header().Set("Content-Type", contentType)
	w.Write(imageData)
}
