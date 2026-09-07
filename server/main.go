package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ChristopherGoodale/habit_task_tracker/server/internal/db"
	"github.com/ChristopherGoodale/habit_task_tracker/server/internal/middleware"
	"github.com/ChristopherGoodale/habit_task_tracker/server/internal/task"
)

func main() {
	databaseURL := envOrDefault("DATABASE_URL", "postgres://habit_tracker:habit_tracker@localhost:5432/habit_tracker?sslmode=disable")
	port := envOrDefault("PORT", "8080")
	clientOrigin := envOrDefault("CLIENT_ORIGIN", "http://localhost:5173")

	conn, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer conn.Close()

	handlers := task.NewHandlers(task.NewStore(conn))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/tasks", handlers.List)
	mux.HandleFunc("POST /api/tasks", handlers.Create)
	mux.HandleFunc("PUT /api/tasks/{id}", handlers.Update)
	mux.HandleFunc("DELETE /api/tasks/{id}", handlers.Delete)

	var h http.Handler = mux
	h = middleware.CORS(clientOrigin)(h)
	h = middleware.Logging(h)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
