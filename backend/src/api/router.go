package api

import (
	"database/sql"
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"

	"task-management-app/api/handlers"
	"task-management-app/api/middleware"
	"task-management-app/config"
)

// SetupRouter configures and returns the HTTP router
func SetupRouter(db *sql.DB, cfg *config.Config, staticFiles embed.FS) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)

	// Session store
	store := sessions.NewCookieStore([]byte(cfg.SessionKey))

	// Static files - serve from root
	staticFS, _ := fs.Sub(staticFiles, "static")
	fileServer := http.FileServer(http.FS(staticFS))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	// API routes
	apiRouter := chi.NewRouter()

	// Auth handler
	authHandler := handlers.NewAuthHandler(db, store)

	// Health check handler (public)
	healthHandler := handlers.NewHealthHandler(db)
	apiRouter.Get("/health", healthHandler.GetHealth)

	// Public routes
	apiRouter.Post("/login", authHandler.Login)

	// Protected routes (require authentication)
	apiRouter.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(store))
		r.Post("/logout", authHandler.Logout)
		r.Get("/me", authHandler.GetCurrentUser)
		
		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AdminMiddleware)
			userHandler := handlers.NewUserHandler(db)
			r.Post("/users", userHandler.CreateUser)
		})
		
		// Task routes (authenticated users)
		taskHandler := handlers.NewTaskHandler(db)
		r.Get("/tasks", taskHandler.GetTasks)
		r.Post("/tasks", taskHandler.CreateTask)
		r.Get("/tasks/{taskId}", taskHandler.GetTask)
		r.Put("/tasks/{taskId}", taskHandler.UpdateTask)
		r.Delete("/tasks/{taskId}", taskHandler.DeleteTask)

		// Dependency routes (authenticated users)
		dependencyHandler := handlers.NewDependencyHandler(db)
		r.Get("/tasks/{taskId}/dependencies", dependencyHandler.GetDependencies)
		r.Post("/tasks/{taskId}/dependencies", dependencyHandler.AddDependency)
		r.Delete("/tasks/{taskId}/dependencies/{dependsOnTaskId}", dependencyHandler.RemoveDependency)
	})

	// TODO: Add more API routes here as they are implemented
	r.Mount("/api", apiRouter)

	return r
}
