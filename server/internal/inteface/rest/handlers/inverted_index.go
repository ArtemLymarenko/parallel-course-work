package handlers

import (
	"fmt"
	tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"
	"parallel-course-work/server/internal/inteface/rest/dto"
)

type InvertedIndexService interface {
	Search(query string) []string
	AddFile(filePath string) error
	RemoveFile(filePath string) error
	GetFileContent(filePath string) ([]byte, error)
}

type Logger interface {
	Log(...interface{})
}

type InvertedIndex struct {
	invIndexService InvertedIndexService
	logger          Logger
}

func NewInvertedIndex(invIndexService InvertedIndexService, logger Logger) *InvertedIndex {
	return &InvertedIndex{
		invIndexService: invIndexService,
		logger:          logger,
	}
}

func (i *InvertedIndex) Search(ctx *tcpRouter.RequestContext) error {
	const op = "InvertedIndex.Search"

	var body dto.SearchRequest
	err := ctx.ShouldParseBodyJSON(&body)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not parse request body",
		}
		//log.Printf("%v: error parsing request body: %v\n", op, err)
		msg := fmt.Sprintf("%v: error parsing request body: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusBadRequest, errorResponse)
	}

	result := i.invIndexService.Search(body.Query)
	response := dto.SearchResponse{
		Files: result,
	}

	return ctx.ResponseJSON(tcpRouter.StatusOK, response)
}

func (i *InvertedIndex) AddFile(ctx *tcpRouter.RequestContext) error {
	const op = "InvertedIndex.AddFile"

	var body dto.AddFileRequest
	err := ctx.ShouldParseBodyJSON(&body)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not parse request body",
		}
		//log.Printf("%v: error parsing request body: %v\n", op, err)
		msg := fmt.Sprintf("%v: error parsing request body: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusBadRequest, errorResponse)
	}

	err = i.invIndexService.AddFile(body.FileName)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not add file",
		}
		//log.Printf("%v: error adding file: %v\n", op, err)
		msg := fmt.Sprintf("%v: error adding file: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusInternalServerError, errorResponse)
	}

	return ctx.ResponseJSON(tcpRouter.StatusOK, nil)
}

func (i *InvertedIndex) GetFileContent(ctx *tcpRouter.RequestContext) error {
	const op = "InvertedIndex.GetFileContent"

	var body dto.GetFileRequest
	err := ctx.ShouldParseBodyJSON(&body)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not parse request body",
		}
		//log.Printf("%v: error parsing request body: %v\n", op, err)
		msg := fmt.Sprintf("%v: error parsing request body: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusBadRequest, errorResponse)
	}

	content, err := i.invIndexService.GetFileContent(body.FileName)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not find the file",
		}
		//log.Printf("%v: error finding the file: %v\n", op, err)
		msg := fmt.Sprintf("%v: error finding the file: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusNotFound, errorResponse)
	}

	response := dto.GetFileResponse{
		FileContent: string(content),
	}
	return ctx.ResponseJSON(tcpRouter.StatusOK, response)
}

func (i *InvertedIndex) RemoveFile(ctx *tcpRouter.RequestContext) error {
	const op = "InvertedIndex.RemoveFile"

	var body dto.RemoveFileRequest
	err := ctx.ShouldParseBodyJSON(&body)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not parse request body",
		}
		//log.Printf("%v: error parsing request body: %v\n", op, err)
		msg := fmt.Sprintf("%v: error parsing request body: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusBadRequest, errorResponse)
	}

	err = i.invIndexService.RemoveFile(body.FileName)
	if err != nil {
		errorResponse := dto.ErrorResponse{
			Message: "could not remove the file",
		}
		//log.Printf("%v: error removing the file: %v\n", op, err)
		msg := fmt.Sprintf("%v: error removing the file: %v", op, err)
		i.logger.Log(msg)
		return ctx.ResponseJSON(tcpRouter.StatusNotFound, errorResponse)
	}

	return ctx.ResponseJSON(tcpRouter.StatusOK, nil)
}
