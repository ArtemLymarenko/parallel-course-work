package handlers

import (
	"encoding/json"
	tcpClient "golang/tcp_client"
	"net/http"
)

type SearchRequestDto struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Files []string `json:"files"`
}

func (h *Handlers) Search() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.FormValue("search")
		searchAnyWord := r.FormValue("search-mode") == "on"

		reqPath := "/index/search"
		if searchAnyWord {
			reqPath = "/index/search-any"
		}

		req := &tcpClient.Request{
			RequestMeta: tcpClient.RequestMeta{
				Path:   reqPath,
				Method: "GET",
			},
			Body: SearchRequestDto{
				Query: query,
			},
		}

		data, err := tcpClient.Fetch(req, 8080, h.env)
		var response tcpClient.Response
		err = json.Unmarshal(data, &response)
		if err != nil {
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

		_ = h.tmpl.Render(w, "history-item", toRender)
	}
}
