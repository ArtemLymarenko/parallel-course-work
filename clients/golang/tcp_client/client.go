package tcpClient

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"parallel-course-work/pkg/streamer"
	"strconv"
)

type Client struct {
	conn net.Conn
}

func New(port int) *Client {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conn: conn,
	}
}

func (c *Client) SendRequest(request *Request) (*Response, error) {
	requestBin, err := request.MarshalJSONBinary()
	if err != nil {
		return nil, errors.New("failed to marshal JSON encode request")
	}

	err = streamer.WriteBuff(c.conn, 2048, requestBin)
	if err != nil {
		return nil, errors.New("failed to send request to server")
	}

	result, err := streamer.ReadBuff(c.conn)
	if err != nil {
		return nil, errors.New("failed to retrieve result from server")
	}

	var response Response
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, errors.New("failed to unmarshal response")
	}

	return &response, nil
}

func (c *Client) CloseConn() {
	_ = c.conn.Close()
}
