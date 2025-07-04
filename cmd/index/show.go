package index

import (
	"myproject/temp/config"

	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/login", nil)
}

func Invalidlogin(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/invalidlogin", nil)
}

func Checker(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/checker", nil)
}

func Success(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/success", nil)

}

func Invalid(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/invalid", nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/logout", nil)
}

func AddPoster(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/addposter", nil)
}

func Sample(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/sample", nil)
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	config.TPL.ExecuteTemplate(w, "/welcome", nil)

}
