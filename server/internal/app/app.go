package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"parallel-course-work/pkg/threadpool"
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

type HandleFunc func() error

type FileIdxHandler struct {
	Method HandleFunc
}

type Router struct {
	FileIdxHandler FileIdxHandler
}

type Server struct {
	port       int
	conn       net.Listener
	threadPool ThreadPool
	router     Router
}

func New(port int, threadPool ThreadPool) *Server {
	return &Server{port: port, threadPool: threadPool}
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
	fmt.Println("server started on port:", s.port)

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

		task := threadpool.NewTask(idx.Add(1), func() error {
			_ = s.ReadMessage(conn)

			//Router
			// Request -> Handler -> Response
			// 1. Parse Request and extract data to proper structure
			// 2. Pass the structure to the handler
			// 3. Handler processes the data and returns the response
			return nil
		})
		s.threadPool.AddTask(task)
	}
}

func (s *Server) ReadMessage(acceptedConn net.Conn) []byte {
	defer acceptedConn.Close()

	buff := make([]byte, 2048)
	for {
		n, err := acceptedConn.Read(buff)
		if err != nil {
			fmt.Printf("error reading: %v\n", err)
			continue
		}
		fmt.Println(string(buff[:n]))
	}

	return []byte{'o', 'k'}
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

	fmt.Println("server stopped")
	wg.Done()
}
