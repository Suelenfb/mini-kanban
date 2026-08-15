package main

import (
	"log"
	"net/http"

	"kanban-backend/handlers"
	"kanban-backend/storage"
)

// corsMiddleware libera o acesso da aplicação React (rodando em outra
// porta durante o desenvolvimento) à API.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	store := storage.NewTaskStore()
	taskHandler := handlers.NewTaskHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/tasks", taskHandler.TasksCollection)
	mux.HandleFunc("/api/tasks/reorder", taskHandler.Reorder)
	mux.HandleFunc("/api/tasks/", taskHandler.TaskItem)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	handler := corsMiddleware(mux)

	addr := ":8080"
	log.Printf("Servidor Kanban rodando em http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
