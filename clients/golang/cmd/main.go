package main

import (
	"golang/app"
	"golang/handlers"
	htmlRender "golang/html_render"
	"log"
	"net/http"
	"os"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

func main() {
	env := app.Env(os.Getenv("ENV"))

	mux := http.NewServeMux()

	tmpl := htmlRender.NewTemplates()

	fs := http.FileServer(http.Dir("views/styles"))
	mux.Handle("/styles/*", http.StripPrefix("/styles/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Render(w, "index", map[string]interface{}{})
	})

	mux.HandleFunc("/search", handlers.Search(tmpl, env))
	mux.HandleFunc("/download", handlers.Download(tmpl, env))
	mux.HandleFunc("/add-file", handlers.AddFile(tmpl, env))
	mux.HandleFunc("/remove-file", handlers.RemoveFile(tmpl, env))

	handler := Logging(mux)

	log.Println("Server listening on :3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", handler))
}
