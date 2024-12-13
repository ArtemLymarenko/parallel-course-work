package tcpRouter

import "errors"

var (
	ErrInvalidRequestMethod = errors.New("invalid request method")
	ErrInvalidRequestPath   = errors.New("invalid request path")
)

type RequestMeta struct {
	Path   RequestPath   `json:"path"`
	Method RequestMethod `json:"method"`
}

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	DELETE RequestMethod = "DELETE"
)

func (requestMethod RequestMethod) Validate() error {
	switch requestMethod {
	case GET, POST, DELETE:
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
