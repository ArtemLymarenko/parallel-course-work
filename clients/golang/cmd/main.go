package main

import (
	"golang/handlers"
	htmlRender "golang/html_render"
	"log"
	"net/http"
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
	mux := http.NewServeMux()

	tmpl := htmlRender.NewTemplates()

	fs := http.FileServer(http.Dir("views/styles"))
	mux.Handle("/styles/*", http.StripPrefix("/styles/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Render(w, "index", map[string]interface{}{})
	})

	mux.HandleFunc("/search", handlers.Search(tmpl))
	mux.HandleFunc("/download", handlers.Download(tmpl))
	mux.HandleFunc("/add-file", handlers.AddFile(tmpl))
	mux.HandleFunc("/remove-file", handlers.RemoveFile(tmpl))

	handler := Logging(mux)

	log.Fatal(http.ListenAndServe("0.0.0.0:3000", handler))
}
