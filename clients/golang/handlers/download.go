package handlers

import (
	"encoding/json"
	tcpClient "golang/tcp_client"
	"net/http"
)

type DownloadRequestDto struct {
	FileName string `json:"fileName"`
}

type GetFileResponse struct {
	FileContent string `json:"fileContent"`
}

func (h *Handlers) Download() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := r.URL.Query().Get("filename")

		req := &tcpClient.Request{
			RequestMeta: tcpClient.RequestMeta{
				Path:   "/index/file",
				Method: "GET",
			},
			Body: DownloadRequestDto{
				FileName: fileName,
			},
		}

		data, err := tcpClient.Fetch(req, 8080, h.env)
		var response tcpClient.Response
		err = json.Unmarshal(data, &response)
		if err != nil || response.Status != tcpClient.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var fileResponse GetFileResponse
		err = json.Unmarshal(response.Body, &fileResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		_ = h.tmpl.Render(w, "file-content", fileResponse)
	}
}
