package tcpRouter

import "errors"

var (
	ErrInvalidRequestMethod = errors.New("invalid request method")
	ErrInvalidRequestPath   = errors.New("invalid request path")
)

type RequestMethod string

const (
	GET  RequestMethod = "GET"
	POST RequestMethod = "POST"
)

func (requestMethod RequestMethod) Validate() error {
	switch requestMethod {
	case GET, POST:
		return nil
	default:
		return ErrInvalidRequestMethod
	}
}

type RequestPath string

func (requestPath RequestPath) Validate() error {
	if requestPath == "" {
		return ErrInvalidRequestPath
	}

	return nil
}
