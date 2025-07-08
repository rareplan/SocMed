package main

import (
	"log"
	"myproject/cmd/index"
	"myproject/entities"

	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")

	log.Println("Server Started...")
	log.Println("Current Port: " + port)
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public"))))
	mux.Handle("/templates/", http.StripPrefix("/templates", http.FileServer(http.Dir("templates"))))

	mux.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("css"))))
	mux.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("js"))))
	mux.Handle("/json-browse/", http.StripPrefix("/json-browse", http.FileServer(http.Dir("json-browse"))))

	//////////// para sa HTML DISPLAY //////////////////////////////////////

	mux.HandleFunc("/login", index.Login)
	mux.HandleFunc("/sample", index.Sample)
	mux.HandleFunc("/invalidlogin", index.Invalidlogin)
	mux.HandleFunc("/welcome", entities.AuthMiddleware([]string{"admin", "checker", "user"}, index.Welcome))
	mux.HandleFunc("/home", entities.AuthMiddleware([]string{"admin", "checker", "user"}, entities.Dashboard))
	mux.HandleFunc("/calendar", entities.AuthMiddleware([]string{"admin", "checker", "user"}, entities.Calendar))
	mux.HandleFunc("/note", entities.AuthMiddleware([]string{"admin", "checker", "user"}, entities.Note))
	mux.HandleFunc("/activity", entities.AuthMiddleware([]string{"admin", "checker", "user"}, entities.Act))
	mux.HandleFunc("/logout", entities.AuthMiddleware([]string{"admin", "checker", "user"}, index.Logout))
	mux.HandleFunc("/invalid", entities.AuthMiddleware([]string{"admin", "checker"}, index.Invalid))
	mux.HandleFunc("/success", entities.AuthMiddleware([]string{"admin", "checker"}, index.Success))
	mux.HandleFunc("/addposter", entities.AuthMiddleware([]string{"admin", "user"}, index.AddPoster))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	//////////// Para sa function at database connection //////////////////////////////////////
	mux.HandleFunc("/loginprocess", entities.LoginProcess)
	mux.HandleFunc("/logoutprocess", entities.LogoutProcess)
	mux.HandleFunc("/updatelink", entities.UpdateLink)
	mux.HandleFunc("/allposter", entities.AuthMiddleware([]string{"admin", "user"}, entities.GetPoster))
	mux.HandleFunc("/useraccess", entities.AuthMiddleware([]string{"admin", "user"}, entities.UserAccess))
	mux.HandleFunc("/allposterchecker", entities.AuthMiddleware([]string{"admin", "checker"}, entities.GetPostercChecker))
	mux.HandleFunc("/updateposter", entities.UpdatePoster)
	mux.HandleFunc("/insertposter", entities.InsertPoster)
	mux.HandleFunc("/deleteposter", entities.DeletePoster)
	mux.HandleFunc("/deleteuser", entities.DeleteUser)
	mux.HandleFunc("/insertuser", entities.InsertUser)
	mux.HandleFunc("/updateuser", entities.UpdateUser)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		panic(err)
	}
}
