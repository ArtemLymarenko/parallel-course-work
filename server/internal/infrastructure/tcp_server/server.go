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
	AddTask(task *threadpool.Task) bool
}

type Router interface {
	Handle(request *tcpRouter.RequestContext) error
	ParseRequest(raw []byte) (*tcpRouter.RequestContext, error)
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

func (s *Server) getAddr() string {
	return fmt.Sprintf(":%d", s.port)
}

func (s *Server) Start() {
	conn, err := net.Listen("tcp", s.getAddr())
	if err != nil {
		log.Fatal(err)
	}

	s.conn = conn
	fmt.Println("Server started on port:", s.port)

	wg := sync.WaitGroup{}
	go s.gracefulShutDown(&wg)

	s.threadPool.MustRun()
	s.AcceptConnections()

	wg.Wait()
}

func (s *Server) AcceptConnections() {
	idx := atomic.Int64{}
	idx.Store(0)
	for {
		conn, err := s.conn.Accept()
		if err != nil {
			log.Fatal(err)
		}

		err = s.HandleConnection(conn, idx.Add(1))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (s *Server) HandleConnection(conn net.Conn, connIdx int64) error {
	rawRequest, err := s.ReadMessage(conn)
	if err != nil {
		return err
	}

	request, err := s.router.ParseRequest(rawRequest)
	if err != nil {
		return err
	}

	request.BindConn(conn)

	task := threadpool.NewTask(connIdx, func() error {
		return s.router.Handle(request)
	})

	ok := s.threadPool.AddTask(task)
	if !ok {
		return fmt.Errorf("task with idx %d was not added", connIdx)
	}

	return nil
}

func (s *Server) ReadMessage(conn net.Conn) ([]byte, error) {
	defer conn.Close()

	var buffer bytes.Buffer
	for {
		chunk := make([]byte, 2048)
		n, err := conn.Read(chunk)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("error reading: %v", err)
		}
		buffer.Write(chunk[:n])
	}

	return buffer.Bytes(), nil
}

func (s *Server) gracefulShutDown(wg *sync.WaitGroup) {
	wg.Add(1)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.conn.Close(); err != nil {
		log.Fatal(err)
	}

	s.threadPool.MustTerminate()

	fmt.Println("Server stopped")
	wg.Done()
}
