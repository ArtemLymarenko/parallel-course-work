package handlers

import (
	"encoding/json"
	htmlRender "golang/html_render"
	tcpClient "golang/tcp_client"
	"net/http"
)

type DownloadRequestDto struct {
	FileName string `json:"fileName"`
}

type GetFileResponse struct {
	FileContent string `json:"fileContent"`
}

func Download(tmpl *htmlRender.Templates) func(w http.ResponseWriter, r *http.Request) {
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

		data, err := tcpClient.Fetch(req, 8080)
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

		tmpl.Render(w, "file-content", fileResponse)
	}
}
