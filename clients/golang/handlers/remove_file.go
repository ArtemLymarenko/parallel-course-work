package handlers

import (
	"encoding/json"
	"net/http"
	htmlRender "parallel-course-work/clients/golang/html_render"
	tcpClient "parallel-course-work/clients/golang/tcp_client"
)

type ErrResponse struct {
	Message string `json:"message,omitempty"`
}

func RemoveFile(tmpl *htmlRender.Templates) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := r.FormValue("file-path")

		req := &tcpClient.Request{
			RequestMeta: tcpClient.RequestMeta{
				Path:   "/index/file",
				Method: "DELETE",
			},
			Body: DownloadRequestDto{
				FileName: filePath,
			},
		}

		data, err := tcpClient.Fetch(req, 8080)
		var response tcpClient.Response
		err = json.Unmarshal(data, &response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var removeErrResponse ErrResponse
		err = json.Unmarshal(response.Body, &removeErrResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		render := struct {
			StatusCode    string
			StatusMessage string
		}{
			StatusCode:    response.Status.String(),
			StatusMessage: removeErrResponse.Message,
		}

		tmpl.Render(w, "status", render)
	}
}
