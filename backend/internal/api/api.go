package api

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Router    *mux.Router
	Templates *template.Template
}

func NewServer() *Server {
	// Create a new router
	router := mux.NewRouter()

	// Parse all templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	return &Server{
		Router:    router,
		Templates: templates,
	}
}

func (s *Server) SetupRoutes() {
	// Serve static files
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	s.Router.PathPrefix("/static-frontend/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static-frontend/_app/") {
			http.StripPrefix("/static-frontend/", http.FileServer(http.Dir("static-frontend/"))).ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, "static-frontend/index.html")
	})

	// Page routers - only GET handlers for now
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.Router.HandleFunc("/login", s.handleLoginPage).Methods("GET")
	s.Router.HandleFunc("/register", s.handleRegisterPage).Methods("GET")
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":    "OurChat - Home",
		"LoggedIn": false,   // No login functionality yet
		"Page":     "index", // Tell the layout which content to include
	}
	s.renderTemplate(w, "layout.html", data)
}

// handleRegisterPage serves the registration page
func (s *Server) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "OurChat - Register",
		"Page":  "register", // Tell the layout which content to include
	}
	s.renderTemplate(w, "layout.html", data)
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "OurChat - Login",
		"Page":  "login", // Tell the layout which content to include
	}

	s.renderTemplate(w, "layout.html", data)
}

// renderTemplate renders a template with the given data
func (s *Server) renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
	// Set default content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Add current year to data for footer
	data["CurrentYear"] = time.Now().Year()

	// Execute the template
	err := s.Templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Error rendering template %s: %v", tmpl, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
