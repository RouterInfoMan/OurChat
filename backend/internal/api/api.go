package api

import (
	"OurChat/internal/db"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Router    *mux.Router
	Templates *template.Template
	DB        *db.DB
}

// TODO After creating the configs
// Load paths and envs from config file

func NewServer() *Server {
	// Create a new router
	router := mux.NewRouter()

	// Parse all templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	database, err := db.NewDB("./data/ourchat.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return &Server{
		Router:    router,
		Templates: templates,
		DB:        database,
	}
}

func (s *Server) SetupRoutes() {
	s.Router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Construct the file path
		filesRoot := "static-frontend"

		if r.URL.Path == "/index" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// Check if the file exists
		if _, err := http.Dir(filesRoot).Open(r.URL.Path); err != nil {
			// If file not found, check for {file}.html
			if _, err := http.Dir(filesRoot).Open(r.URL.Path + ".html"); err != nil {
				// If file not found, serve 200.html
				http.ServeFile(w, r, filesRoot+"/200.html")
				return
			}

			http.ServeFile(w, r, filesRoot+r.URL.Path+".html")
			return
		}

		// Serve the requested file
		http.FileServer(http.Dir(filesRoot)).ServeHTTP(w, r)
	})
}

// func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]interface{}{
// 		"Title":    "OurChat - Home",
// 		"LoggedIn": false,   // No login functionality yet
// 		"Page":     "index", // Tell the layout which content to include
// 	}
// 	s.renderTemplate(w, "layout.html", data)
// }

// // handleRegisterPage serves the registration page
// func (s *Server) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]interface{}{
// 		"Title": "OurChat - Register",
// 		"Page":  "register", // Tell the layout which content to include
// 	}
// 	s.renderTemplate(w, "layout.html", data)
// }

// func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
// 	data := map[string]interface{}{
// 		"Title": "OurChat - Login",
// 		"Page":  "login", // Tell the layout which content to include
// 	}

// 	s.renderTemplate(w, "layout.html", data)
// }

// // renderTemplate renders a template with the given data
// func (s *Server) renderTemplate(w http.ResponseWriter, tmpl string, data map[string]interface{}) {
// 	// Set default content type
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")

// 	// Add current year to data for footer
// 	data["CurrentYear"] = time.Now().Year()

// 	// Execute the template
// 	err := s.Templates.ExecuteTemplate(w, tmpl, data)
// 	if err != nil {
// 		log.Printf("Error rendering template %s: %v", tmpl, err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 	}
// }
