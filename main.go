package main

import (
	"html/template"
	"log"
	"net/http"

	"multi-society-portal/database"

	"github.com/joho/godotenv"
)

type LoginPageData struct {
	Society string
	Success string
}

func selectSocietyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/select_society.html"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	selectedSociety := r.URL.Query().Get("society")
	successMsg := r.URL.Query().Get("success")

	data := LoginPageData{
		Society: selectedSociety,
		Success: successMsg,
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

func committeeRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	selectedSociety := r.URL.Query().Get("society")

	data := LoginPageData{
		Society: selectedSociety,
	}

	tmpl := template.Must(template.ParseFiles("templates/committee_registration.html"))
	tmpl.Execute(w, data)
}
func registerResidentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fullName := r.FormValue("full_name")
	flatNumber := r.FormValue("flat_number")
	mobile := r.FormValue("mobile")
	email := r.FormValue("email")
	wing := r.FormValue("wing")

	_, err := database.DB.Exec(
		`INSERT INTO residents
    (full_name, flat_number, mobile, email, wing)
    VALUES ($1, $2, $3, $4, $5)`,
		fullName,
		flatNumber,
		mobile,
		email,
		wing,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success redirect
	http.Redirect(
		w,
		r,
		"/login?success=Resident Registration Successful",
		http.StatusSeeOther,
	)
}

func committeeRegisterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fullName := r.FormValue("full_name")
	email := r.FormValue("email")
	mobile := r.FormValue("mobile")
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")
	designation := r.FormValue("designation")

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(
		`INSERT INTO committee_members
    (full_name, email, mobile, username, designation, password_hash)
    VALUES ($1, $2, $3, $4, $5, $6)`,
		fullName,
		email,
		mobile,
		username,
		designation,
		password,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Success redirect
	http.Redirect(
		w,
		r,
		"/login?success=Committee Registration Successful",
		http.StatusSeeOther,
	)
}

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect PostgreSQL
	database.ConnectDB()

	// Static files
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	// Routes
	http.HandleFunc("/", selectSocietyHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/registration", registrationHandler)
	http.HandleFunc("/committee-registration", committeeRegistrationHandler)
	http.HandleFunc("/register", registerResidentHandler)
	http.HandleFunc("/committee-register", committeeRegisterHandler)
	log.Println("🚀 Server started at :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
