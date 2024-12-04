package tcpServer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"parallel-course-work/pkg/threadpool"
	tcpRouter "parallel-course-work/server/internal/infrastructure/tcp_server/router"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type ThreadPool interface {
	MustRun()
	MustTerminate()
	AddTask(task *threadpool.Task) error
}

type Router interface {
	Handle(request *tcpRouter.RequestContext) error
	ParseRequest(raw []byte) (*tcpRouter.Request, error)
}

type Server struct {
	port       int
	conn       net.Listener
	threadPool ThreadPool
	router     Router
}

func New(port int, threadPool ThreadPool, router Router) *Server {
	return &Server{
		port:       port,
		threadPool: threadPool,
		router:     router,
	}
}

func (server *Server) getAddr() string {
	return fmt.Sprintf(":%d", server.port)
}

func (server *Server) Start() {
	conn, err := net.Listen("tcp", server.getAddr())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server.conn = conn
	fmt.Println("Server started on port:", server.port)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go server.gracefulShutDown(&wg)

	server.threadPool.MustRun()
	server.acceptConnections()

	wg.Wait()
}

func (server *Server) acceptConnections() {
	idx := atomic.Int64{}
	idx.Store(0)
	for {
		conn, err := server.conn.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		fmt.Println("New connection accepted")
		err = server.handleConnection(conn, idx.Add(1))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (server *Server) handleConnection(clientConn net.Conn, connIdx int64) error {
	rawRequest, err := server.readMessage(clientConn)
	if err != nil {
		return err
	}
	fmt.Println("Received message:", rawRequest)

	request, err := server.router.ParseRequest(rawRequest)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return err
	}
	fmt.Println("Parsed:", request)
	requestCtx := tcpRouter.NewRequestContext(request, clientConn)

	task := threadpool.NewTask(connIdx, func() error {
		return server.router.Handle(requestCtx)
	})

	return server.threadPool.AddTask(task)
}

func (server *Server) readMessage(clientConn net.Conn) ([]byte, error) {
	lengthBuffer := make([]byte, 4)
	_, err := clientConn.Read(lengthBuffer)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %v", err)
	}
	messageLength := int(lengthBuffer[0])<<24 | int(lengthBuffer[1])<<16 | int(lengthBuffer[2])<<8 | int(lengthBuffer[3])

	var (
		buffer      bytes.Buffer
		totalLength int
	)

	for {
		chunk := make([]byte, 2048)
		n, err := clientConn.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading: %v", err)
		}

		buffer.Write(chunk[:n])
		totalLength += n
		if totalLength == messageLength {
			break
		}
	}

	return buffer.Bytes(), nil
}

func (server *Server) gracefulShutDown(wg *sync.WaitGroup) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.conn.Close(); err != nil {
		log.Fatal(err)
	}

	server.threadPool.MustTerminate()

	fmt.Println("Server stopped")
	wg.Done()
}
