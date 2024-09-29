package main

import (
	"fmt"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"log"
	"net/http"
	"time"
)

func setupRoutes(mux *http.ServeMux, buff tasks.Buffer) {
	mux.HandleFunc("/tasks", TaskRouter(buff))
	mux.HandleFunc("/tasks/", TaskRouter(buff))
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("%s %s %v\n", r.Method, r.URL, time.Since(start))
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Sqlite3 file
	sqlFileName := "./tmp/tasks.db"
	port := "8080"

	buff, err := tasks.NewSqlite3Buffer(sqlFileName)

	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	setupRoutes(mux, buff)

	handler := logMiddleware(mux)
	handler = panicMiddleware(handler)

	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
