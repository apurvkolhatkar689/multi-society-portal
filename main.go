package main

import (
	"html/template"
	"log"
	"net/http"
)

type LoginPageData struct {
	Society string
}

func selectSocietyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/select_society.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	selectedSociety := r.URL.Query().Get("society")
	data := LoginPageData{
		Society: selectedSociety,
	}
	tmpl := template.Must(template.ParseFiles("templates/login.html"))
	tmpl.Execute(w, data)
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	selectedSociety := r.URL.Query().Get("society")
	data := LoginPageData{
		Society: selectedSociety,
	}
	tmpl := template.Must(template.ParseFiles("templates/registration.html"))
	tmpl.Execute(w, data)
}

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", selectSocietyHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/registration", registrationHandler)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
