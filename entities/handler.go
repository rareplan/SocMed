package entities

import (
	"bytes"
	"myproject/temp/config"
	"net/http"
)

// ///////// LOGIN HANDLER //////////
func LoginProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	loginHandler(w, r)

}

// //// PARA SA LOGOUT PROCESS //////
func LogoutProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	logoutProcessHandler(w, r)

}

// //// PARA SA UPDATE POSTER //////
func UpdatePoster(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	UpdatePosterHandle(w, r)

}

// ///// PARA SA UPDATE LINK //////
func UpdateLink(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	UpdateLinkHandle(w, r)

}

// ///// PARA SA INSERT POSTER //////
func InsertPoster(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	InsertPosterHandle(w, r)

}

// //// PARA SA DELETE POSTER //////
func DeletePoster(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	DeletePosterHandle(w, r)

}

// //// PARA SA GET POSTER //////
func GetPoster(w http.ResponseWriter, r *http.Request) {
	// ✅ 1. Check for session via cookie
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

	// ✅ 2. Get posters from DB/service
	posters, err := AllPoster()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ 3. Combine role flags + posters for template
	data := struct {
		IsAdmin    bool
		IsDesigner bool
		IsVerifier bool
		Role       string
		Posters    any
	}{
		IsAdmin:    session.Role == "admin",
		IsDesigner: session.Role == "designer",
		IsVerifier: session.Role == "verifier",
		Role:       session.Role,
		Posters:    posters,
	}

	var tplOutput bytes.Buffer
	// ✅ 4. Render correct template name
	err = config.TPL.ExecuteTemplate(w, "poster", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
		return
	}
	tplOutput.WriteTo(w)
}

// ////// PARA SA GET POSTER CHECKER //////
func GetPostercChecker(w http.ResponseWriter, r *http.Request) {
	// ✅ 1. Check for session via cookie
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

	// ✅ 2. Get posters from DB/service
	posters, err := AllPoster()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ 3. Combine role flags + posters for template
	data := struct {
		IsAdmin    bool
		IsDesigner bool
		IsVerifier bool
		Role       string
		Posters    any
	}{
		IsAdmin:    session.Role == "admin",
		IsDesigner: session.Role == "designer",
		IsVerifier: session.Role == "verifier",
		Role:       session.Role,
		Posters:    posters,
	}

	// ✅ 4. Render correct template name
	var tplOutput bytes.Buffer
	err = config.TPL.ExecuteTemplate(w, "allchecker", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
		return
	}
	tplOutput.WriteTo(w)
}

// ///// PARA SA USER ACCESS //////
func UserAccess(w http.ResponseWriter, r *http.Request) {
	// ✅ 1. Check for session via cookie
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

	// ✅ 2. Get posters from DB/service
	users, err := AllUsers()
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ 3. Combine role flags + posters for template
	data := struct {
		IsAdmin    bool
		IsDesigner bool
		IsVerifier bool
		Role       string
		Users      any
	}{
		IsAdmin:    session.Role == "admin",
		IsDesigner: session.Role == "designer",
		IsVerifier: session.Role == "verifier",
		Role:       session.Role,
		Users:      users,
	}

	// ✅ 4. Render correct template name
	var tplOutput bytes.Buffer
	err = config.TPL.ExecuteTemplate(w, "user", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
		return
	}
	tplOutput.WriteTo(w)
}

// /// PARA SA DELETE USER //////
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	DeleteUserHandle(w, r)

}

// /// PARA SA INSERT USER //////
func InsertUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	InsertUserHandle(w, r)

}

// //// PARA SA UPDATE USER //////
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	UpdateUserHandle(w, r)

}

func InsertDataUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	InsertData(w, r)

}
