package main

import (
	"log"
	"net/http"
	"parallel-course-work/clients/golang/handlers"
	htmlRender "parallel-course-work/clients/golang/html_render"
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
		tmpl.Render(w, "index", struct{}{})
	})

	mux.HandleFunc("/search", handlers.Search(tmpl))
	mux.HandleFunc("/download", handlers.Download(tmpl))

	handler := Logging(mux)

	log.Fatal(http.ListenAndServe(":3000", handler))
}
