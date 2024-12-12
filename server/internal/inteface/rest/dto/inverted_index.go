package dto

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Files []string `json:"files"`
}
