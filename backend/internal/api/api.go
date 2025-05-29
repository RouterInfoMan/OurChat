package api

import (
	"log"

	"OurChat/internal/api/handlers"
	"OurChat/internal/api/middleware"
	"OurChat/internal/db"

	"github.com/gorilla/mux"
)

type Server struct {
	Router         *mux.Router
	DB             *db.DB
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	ChatHandler    *handlers.ChatHandler
	MessageHandler *handlers.MessageHandler
	MediaHandler   *handlers.MediaHandler
	AuthMiddleware *middleware.AuthMiddleware
}

// NewServer creates a new API server
func NewServer() *Server {
	// Create a new router
	router := mux.NewRouter()

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
	mediaHandler := handlers.NewMediaHandler(database)

	// Create middleware
	authMiddleware := middleware.NewAuthMiddleware(database)

	return &Server{
		Router:         router,
		DB:             database,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		ChatHandler:    chatHandler,
		MessageHandler: messageHandler,
		MediaHandler:   mediaHandler,
		AuthMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
	// API Routes TEMP FIX
	api := s.Router.PathPrefix("/api").Subrouter()

	// Auth routes - no authentication required
	api.HandleFunc("/register", s.AuthHandler.HandleRegister).Methods("POST")
	api.HandleFunc("/login", s.AuthHandler.HandleLogin).Methods("POST")
	api.HandleFunc("/request-password-reset", s.AuthHandler.HandleRequestPasswordReset).Methods("POST")
	api.HandleFunc("/reset-password", s.AuthHandler.HandleResetPassword).Methods("POST")

	// Protected routes - authentication required
	protected := api.PathPrefix("").Subrouter()
	protected.Use(s.AuthMiddleware.Middleware)

	// User routes
	protected.HandleFunc("/logout", s.AuthHandler.HandleLogout).Methods("POST")
	protected.HandleFunc("/profile", s.UserHandler.HandleGetProfile).Methods("GET")
	protected.HandleFunc("/profile", s.UserHandler.HandleUpdateProfile).Methods("PUT")
	protected.HandleFunc("/profile/picture", s.MediaHandler.HandleUploadProfilePicture).Methods("POST")

	// Media routes
	protected.HandleFunc("/media/upload", s.MediaHandler.HandleUploadMedia).Methods("POST")
	protected.HandleFunc("/media/{type}/{filename}", s.MediaHandler.HandleServeMedia).Methods("GET")

	// Chat routes
	protected.HandleFunc("/chats", s.ChatHandler.HandleGetChats).Methods("GET")
	protected.HandleFunc("/chats", s.ChatHandler.HandleCreateChat).Methods("POST")
	protected.HandleFunc("/chats/{chatID}", s.ChatHandler.HandleGetChat).Methods("GET")
	protected.HandleFunc("/chats/{chatID}/members", s.ChatHandler.HandleGetChatMembers).Methods("GET")

	// Message routes
	protected.HandleFunc("/chats/{chatID}/messages", s.MessageHandler.HandleGetMessages).Methods("GET")
	protected.HandleFunc("/chats/{chatID}/messages", s.MessageHandler.HandleSendMessage).Methods("POST")
	protected.HandleFunc("/chats/{chatID}/messages/read", s.MessageHandler.HandleMarkMessagesAsRead).Methods("POST")
	protected.HandleFunc("/chats/{chatID}/messages/search", s.MessageHandler.HandleSearchMessages).Methods("GET")
	protected.HandleFunc("/chats/{chatID}/messages/media", s.MessageHandler.HandleSendMediaMessage).Methods("POST")

	// Helper routes
	protected.HandleFunc("/users/search", s.UserHandler.HandleSearchUsers).Methods("GET")
	protected.HandleFunc("/users", s.UserHandler.HandleGetUsersByIDs).Methods("GET", "POST")

	// s.Router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// Construct the file path
	// 	filesRoot := "static-frontend"

	// 	if r.URL.Path == "/index" {
	// 		http.Redirect(w, r, "/", http.StatusFound)
	// 		return
	// 	}

	// 	// Check if the file exists
	// 	if _, err := http.Dir(filesRoot).Open(r.URL.Path); err != nil {
	// 		// If file not found, check for {file}.html
	// 		if _, err := http.Dir(filesRoot).Open(r.URL.Path + ".html"); err != nil {
	// 			// If file not found, serve 200.html
	// 			http.ServeFile(w, r, filesRoot+"/200.html")
	// 			return
	// 		}

	// 		http.ServeFile(w, r, filesRoot+r.URL.Path+".html")
	// 		return
	// 	}

	// 	// Serve the requested file
	// 	http.FileServer(http.Dir(filesRoot)).ServeHTTP(w, r)
	// })

}
