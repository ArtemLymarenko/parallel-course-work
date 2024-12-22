package tcpClient

import (
	"errors"
	"net"
	"parallel-course-work/pkg/streamer"
	"strconv"
)

func Fetch(request *Request, port int) ([]byte, error) {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
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
