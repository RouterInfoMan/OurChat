package api

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"OurChat/internal/api/handlers"
	"OurChat/internal/api/middleware"
	"OurChat/internal/db"

	"github.com/gorilla/mux"
)

type Server struct {
	Router         *mux.Router
	Templates      *template.Template
	DB             *db.DB
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	ChatHandler    *handlers.ChatHandler
	MessageHandler *handlers.MessageHandler
	AuthMiddleware *middleware.AuthMiddleware
}

// NewServer creates a new API server
func NewServer() *Server {
	// Create a new router
	router := mux.NewRouter()

	// Parse all templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Connect to database
	database, err := db.NewDB("./data/ourchat.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Create handlers
	authHandler := handlers.NewAuthHandler(database)
	userHandler := handlers.NewUserHandler(database)
	chatHandler := handlers.NewChatHandler(database)
	messageHandler := handlers.NewMessageHandler(database)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(database)

	return &Server{
		Router:         router,
		Templates:      templates,
		DB:             database,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ChatHandler:    chatHandler,
		MessageHandler: messageHandler,
		AuthMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all the routes for the server
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

	// Page routes - only GET handlers for now
	s.Router.HandleFunc("/", s.handleIndex).Methods("GET")
	s.Router.HandleFunc("/login", s.handleLoginPage).Methods("GET")
	s.Router.HandleFunc("/register", s.handleRegisterPage).Methods("GET")

	// API Routes
	api := s.Router.PathPrefix("/api").Subrouter()

	// Auth routes - no authentication required
	api.HandleFunc("/register", s.AuthHandler.HandleRegister).Methods("POST")
	api.HandleFunc("/login", s.AuthHandler.HandleLogin).Methods("POST")

	// Protected routes - authentication required
	protected := api.PathPrefix("").Subrouter()
	protected.Use(s.AuthMiddleware.Middleware)

	// User routes
	protected.HandleFunc("/logout", s.AuthHandler.HandleLogout).Methods("POST")
	protected.HandleFunc("/profile", s.UserHandler.HandleGetProfile).Methods("GET")
	protected.HandleFunc("/profile", s.UserHandler.HandleUpdateProfile).Methods("PUT")

	// Chat routes
	protected.HandleFunc("/chats", s.ChatHandler.HandleGetChats).Methods("GET")
	protected.HandleFunc("/chats", s.ChatHandler.HandleCreateChat).Methods("POST")
	protected.HandleFunc("/chats/{chatID}", s.ChatHandler.HandleGetChat).Methods("GET")

	// Message routes
	protected.HandleFunc("/chats/{chatID}/messages", s.MessageHandler.HandleGetMessages).Methods("GET")
	protected.HandleFunc("/chats/{chatID}/messages", s.MessageHandler.HandleSendMessage).Methods("POST")
	protected.HandleFunc("/chats/{chatID}/messages/read", s.MessageHandler.HandleMarkMessagesAsRead).Methods("POST")
	protected.HandleFunc("/chats/{chatID}/messages/search", s.MessageHandler.HandleSearchMessages).Methods("GET")
}

// Page template handlers

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":    "OurChat - Home",
		"LoggedIn": false,
		"Page":     "index",
	}
	s.renderTemplate(w, "layout.html", data)
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "OurChat - Login",
		"Page":  "login",
	}
	s.renderTemplate(w, "layout.html", data)
}

func (s *Server) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "OurChat - Register",
		"Page":  "register",
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
