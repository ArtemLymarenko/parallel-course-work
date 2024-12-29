package handlers

import (
	"encoding/json"
	"golang/app"
	htmlRender "golang/html_render"
	tcpClient "golang/tcp_client"
	"net/http"
)

func AddFile(tmpl *htmlRender.Templates, env app.Env) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := r.FormValue("file-path")

		req := &tcpClient.Request{
			RequestMeta: tcpClient.RequestMeta{
				Path:   "/index/file",
				Method: "POST",
			},
			Body: DownloadRequestDto{
				FileName: filePath,
			},
		}

		data, err := tcpClient.Fetch(req, 8080, env)
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
