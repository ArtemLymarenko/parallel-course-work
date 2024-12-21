package tcpClient

import "encoding/json"

type RequestMeta struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type Request struct {
	RequestMeta     RequestMeta `json:"meta"`
	Body            any         `json:"body,omitempty"`
	ConnectionAlive bool        `json:"connectionAlive,omitempty"`
}

func (r *Request) MarshalJSONBinary() ([]byte, error) {
	return json.Marshal(r)
}
