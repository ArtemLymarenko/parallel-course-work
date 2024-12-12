package handlers

import (
	"encoding/json"
	"log"
	tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"
	"parallel-course-work/server/internal/inteface/rest/dto"
)

type InvertedIndexService interface {
	Search(query string) []string
}

type InvertedIndex struct {
	invIndexService InvertedIndexService
}

func NewInvertedIndex(invIndexService InvertedIndexService) *InvertedIndex {
	return &InvertedIndex{
		invIndexService: invIndexService,
	}
}

func (i *InvertedIndex) Search(ctx *tcpRouter.RequestContext) error {
	var body dto.SearchRequest
	err := json.Unmarshal(ctx.Request.Body, &body)

	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not parse request body",
		}
		if err = ctx.ResponseJSON(tcpRouter.BadRequest, errorResponse); err != nil {
			log.Println("error parsing request body:", err)
			return err
		}
		return nil
	}

	result := i.invIndexService.Search(body.Query)
	response := dto.SearchResponse{
		Files: result,
	}

	return ctx.ResponseJSON(tcpRouter.OK, response)
}
