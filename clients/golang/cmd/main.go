package main

import (
	"fmt"
	tcpClient "parallel-course-work/clients/golang/tcp_client"
)

type Dto struct {
	Query string `json:"query"`
}

func main() {
	req := &tcpClient.Request{
		RequestMeta: tcpClient.RequestMeta{
			Path:   "/index/search",
			Method: "GET",
		},
		Body: Dto{
			Query: "Hello world",
		},
	}

	client := tcpClient.New(8080)
	defer client.CloseConn()

	response, err := client.SendRequest(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}
