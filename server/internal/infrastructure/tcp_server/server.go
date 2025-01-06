package tcpServer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	tcpRouter "server/internal/infrastructure/tcp_server/router"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ArtemLymarenko/parallel-course-work/pkg/streamer"
	"github.com/ArtemLymarenko/parallel-course-work/pkg/threadpool"
)

type ThreadPool interface {
	MustRun(threadCount int)
	MustTerminate()
	AddTask(task *threadpool.Task) error
}

type Router interface {
	Handle(req *tcpRouter.Request, conn net.Conn) error
	ParseRawRequest(raw []byte) (*tcpRouter.Request, error)
}

const AliveTimeout = 15 * time.Second

type Logger interface {
	Log(...interface{})
}

type Server struct {
	port           int
	conn           net.Listener
	threadPool     ThreadPool
	router         Router
	shutdownSignal chan struct{}
	taskIds        atomic.Int64
	logger         Logger
}

func New(port int, threadPool ThreadPool, router Router, logger Logger) *Server {
	return &Server{
		port:           port,
		threadPool:     threadPool,
		router:         router,
		shutdownSignal: make(chan struct{}),
		taskIds:        atomic.Int64{},
		logger:         logger,
	}
}

func (server *Server) getAddr() string {
	return fmt.Sprintf("0.0.0.0:%d", server.port)
}

func (server *Server) Start(threadsCount int) error {
	conn, err := net.Listen("tcp", server.getAddr())
	if err != nil {
		return err
	}
	defer conn.Close()

	server.conn = conn

	wg := sync.WaitGroup{}
	wg.Add(1)
	go server.gracefulShutDown(&wg)

	server.threadPool.MustRun(threadsCount)
	server.logger.Log("Server started on port:", server.port)

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
			server.logger.Log("Shutting down acceptConnections...")
			return
		default:
			conn, err := server.conn.Accept()
			if err != nil {
				continue
			}

			connIdx := idx.Add(1)
			msg := fmt.Sprintf("client connection [%v] opened", connIdx)
			server.logger.Log(msg)
			if err := server.handleConnections(conn, connIdx); err != nil {
				msg := fmt.Sprintf("error happened %v", err)
				server.logger.Log(msg)
			}
		}
	}
}

func (server *Server) scheduleTask(request *tcpRouter.Request, clientConn net.Conn, connIdx int64) (err error) {
	task := threadpool.NewTask(server.taskIds.Add(1), func() error {
		defer func() {
			if !request.ConnectionAlive {
				if err = clientConn.Close(); err != nil {
					msg := fmt.Sprintf("error occurred: %v", err)
					server.logger.Log(msg)
				}

				msg := fmt.Sprintf("client [%v] disconnected", connIdx)
				server.logger.Log(msg)
			}
		}()

		return server.router.Handle(request, clientConn)
	})

	msg := fmt.Sprintf("Request: method: %v - path: %v", request.RequestMeta.Method, request.RequestMeta.Path)
	server.logger.Log(msg)

	err = server.threadPool.AddTask(task)
	if err != nil {
		return err
	}

	return nil
}

func (server *Server) handleConnections(clientConn net.Conn, connIdx int64) error {
	rawRequest, err := streamer.ReadBuff(clientConn)
	if err != nil {
		return err
	}

	request, err := server.router.ParseRawRequest(rawRequest)
	if err != nil {
		return err
	}

	err = server.scheduleTask(request, clientConn, connIdx)
	if err != nil {
		return err
	}

	if request.ConnectionAlive {
		go func() {
			err = server.handleConnectionAlive(clientConn, connIdx, AliveTimeout)
			if err != nil {
				msg := fmt.Sprintf("error occurred: %v", err)
				server.logger.Log(msg)
			}
		}()
	}

	return nil
}

func (server *Server) handleSingleRequestAlive(clientConn net.Conn, connIdx int64, timeout time.Duration) error {
	if err := clientConn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		msg := fmt.Sprintf("failed to set read deadline: %v", err)
		server.logger.Log(msg)
		return err
	}

	rawRequest, err := streamer.ReadBuff(clientConn)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			msg := fmt.Sprintf("client [%v] timed out", connIdx)
			server.logger.Log(msg)

			_ = clientConn.SetReadDeadline(time.Time{})
			return io.EOF
		}

		return err
	}

	request, err := server.router.ParseRawRequest(rawRequest)
	if err != nil {
		return err
	}

	return server.scheduleTask(request, clientConn, connIdx)
}

func (server *Server) handleConnectionAlive(
	clientConn net.Conn,
	connIdx int64,
	timeout time.Duration,
) error {
	defer clientConn.Close()

	for {
		err := server.handleSingleRequestAlive(clientConn, connIdx, timeout)
		if err != nil {
			if errors.Is(err, io.EOF) {
				msg := fmt.Sprintf("client [%v] disconnected", connIdx)
				server.logger.Log(msg)
				break
			}

			return err
		}
	}

	return nil
}

func (server *Server) gracefulShutDown(wg *sync.WaitGroup) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	close(server.shutdownSignal)
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.conn.Close(); err != nil {
		log.Println(err)
	}

	server.threadPool.MustTerminate()

	server.logger.Log("server stopped")
	wg.Done()
}
