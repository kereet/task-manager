package main

import (
	"log"
	"net/http"
	"os"
	"task-manager/internal/auth"
	"task-manager/internal/database"
	"task-manager/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	log.Printf("Starting server on port %s", serverPort)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("Connected to database")

	taskStore := database.NewTaskStore(db)
	userStore := database.NewUserStore(db)

	taskHandler := handlers.NewHandlers(taskStore)
	authHandler := handlers.NewAuthHandlers(userStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", methodHandler(authHandler.Register, http.MethodPost))
	mux.HandleFunc("/api/login", methodHandler(authHandler.Login, http.MethodPost))

	mux.HandleFunc("/tasks", authMiddleware(methodHandler(taskHandler.GetAllTasks, http.MethodGet)))
	mux.HandleFunc("/tasks/create", authMiddleware(methodHandler(taskHandler.CreateTask, http.MethodPost)))
	mux.HandleFunc("/tasks/", authMiddleware(taskIDHandler(taskHandler)))

	loggedMux := loggingMiddleware(mux)
	serverAddr := ":" + serverPort

	err = http.ListenAndServe(serverAddr, loggedMux)

	if err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", string(rune(claims.UserID)))

		next(w, r)
	}
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func taskIDHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
