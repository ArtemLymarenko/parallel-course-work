package tcpClient

import (
	"errors"
	"fmt"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/streamer"
	"net"
)

func Fetch(request *Request, port int) ([]byte, error) {
	connPath := fmt.Sprintf("server-app:%d", port)

	conn, err := net.Dial("tcp", connPath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	requestBin, err := request.MarshalJSONBinary()
	if err != nil {
		return nil, errors.New("failed to marshal JSON encode request")
	}

	err = streamer.WriteBuff(conn, 2048, requestBin)
	if err != nil {
		return nil, errors.New("failed to send request to server")
	}

	result, err := streamer.ReadBuff(conn)
	if err != nil {
		return nil, errors.New("failed to retrieve result from server")
	}

	return result, nil
}
