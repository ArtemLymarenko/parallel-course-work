package dto

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Files []string `json:"files"`
}

type AddFileRequest struct {
	FileName string `json:"fileName"`
}

type GetFileRequest struct {
	FileName string `json:"fileName"`
}

type GetFileResponse struct {
	FileContent string `json:"fileContent"`
}

type RemoveFileRequest struct {
	FileName string `json:"fileName"`
}
