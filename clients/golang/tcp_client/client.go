package tcpClient

import (
	"errors"
	"fmt"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/streamer"
	"golang/app"
	"net"
)

func GetConnPath(port int, env app.Env) string {
	if env.IsProduction() {
		return fmt.Sprintf("server-app:%d", port)
	}

	return fmt.Sprintf("0.0.0.0:%d", port)
}

func Fetch(request *Request, port int, env app.Env) ([]byte, error) {
	connPath := GetConnPath(port, env)
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
