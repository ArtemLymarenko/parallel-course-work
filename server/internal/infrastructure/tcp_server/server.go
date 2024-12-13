package tcpServer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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
	AddTask(task *threadpool.Task) error
}

type Router interface {
	Handle(raw []byte, conn net.Conn) error
}

type Server struct {
	port           int
	conn           net.Listener
	threadPool     ThreadPool
	router         Router
	shutdownSignal chan struct{}
	taskIds        atomic.Int64
}

func New(port int, threadPool ThreadPool, router Router) *Server {
	return &Server{
		port:           port,
		threadPool:     threadPool,
		router:         router,
		shutdownSignal: make(chan struct{}),
		taskIds:        atomic.Int64{},
	}
}

func (server *Server) getAddr() string {
	return fmt.Sprintf(":%d", server.port)
}

func (server *Server) Start() error {
	conn, err := net.Listen("tcp", server.getAddr())
	if err != nil {
		return err
	}
	defer conn.Close()
	server.conn = conn
	log.Printf("Server started on port: %v\n", server.port)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go server.gracefulShutDown(&wg)

	server.threadPool.MustRun()
	server.acceptConnections()

	wg.Wait()

	return nil
}

func (server *Server) acceptConnections() {
	idx := atomic.Int64{}
	idx.Store(0)

	for {
		select {
		case <-server.shutdownSignal:
			log.Println("Shutting down acceptConnections...")
			return
		default:
			conn, err := server.conn.Accept()
			if err != nil {
				continue
			}

			connIdx := idx.Add(1)
			log.Printf("client connection [%v] opened\n", connIdx)
			go func() {
				err = server.handleConnectionAlive(conn, connIdx)
				if err != nil {
					log.Printf("error occurred: %v\n", err)
				}
			}()
		}
	}
}

func (server *Server) handleConnection(clientConn net.Conn, connIdx int64) error {
	rawRequest, err := server.readMessage(clientConn)
	if err != nil {
		if err == io.EOF {
			log.Printf("client [%v] disconnected\n", connIdx)
			return nil
		}

		return err
	}

	task := threadpool.NewTask(server.taskIds.Add(1), func() error {
		defer func() {
			if err := clientConn.Close(); err != nil {
				log.Printf("error occurred: %v\n", err)
			}
		}()
		return server.router.Handle(rawRequest, clientConn)
	})

	return server.threadPool.AddTask(task)
}

func (server *Server) handleConnectionAlive(clientConn net.Conn, connIdx int64) error {
	defer func() {
		if err := clientConn.Close(); err != nil {
			log.Printf("error occurred: %v\n", err)
		}
	}()

	for {
		rawRequest, err := server.readMessage(clientConn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Printf("client [%v] disconnected\n", connIdx)
				break
			}

			return err
		}

		task := threadpool.NewTask(server.taskIds.Add(1), func() error {
			return server.router.Handle(rawRequest, clientConn)
		})

		err = server.threadPool.AddTask(task)
		if err != nil {
			return err
		}
	}

	return nil
}

func (server *Server) readMessage(clientConn net.Conn) ([]byte, error) {
	lengthBuffer := make([]byte, 4)
	_, err := clientConn.Read(lengthBuffer)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, err
		}
		return nil, errors.New("error reading message length")
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
			if errors.Is(err, io.EOF) {
				return nil, err
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

	close(server.shutdownSignal)
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.conn.Close(); err != nil {
		log.Fatal(err)
	}

	server.threadPool.MustTerminate()

	log.Printf("server stopped\n")
	wg.Done()
}
