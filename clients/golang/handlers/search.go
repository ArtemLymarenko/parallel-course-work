package handlers

import (
	"encoding/json"
	"net/http"
	htmlRender "parallel-course-work/clients/golang/html_render"
	tcpClient "parallel-course-work/clients/golang/tcp_client"
)

type SearchRequestDto struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Files []string `json:"files"`
}

func Search(tmpl *htmlRender.Templates) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("search")

		req := &tcpClient.Request{
			RequestMeta: tcpClient.RequestMeta{
				Path:   "/index/search",
				Method: "GET",
			},
			Body: SearchRequestDto{
				Query: query,
			},
		}

		data, err := tcpClient.Fetch(req, 8080)
		var response tcpClient.Response
		err = json.Unmarshal(data, &response)
		if err != nil || response.Status != tcpClient.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var searchResponse SearchResponse
		err = json.Unmarshal(response.Body, &searchResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		toRender := struct {
			Query string
			Files []string
		}{
			Query: query,
			Files: searchResponse.Files,
		}

		_ = tmpl.Render(w, "history-item", toRender)
	}
}
