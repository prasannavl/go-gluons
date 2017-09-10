package gosock

import (
	"io"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prasannavl/go-gluons/log"
	"github.com/rsms/gotalk"
)

type rwc struct {
	reader io.Reader
	conn   *websocket.Conn
}

func (c *rwc) Write(p []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *rwc) Read(p []byte) (int, error) {
	for {
		if c.reader == nil {
			// Advance to next message.
			var err error
			_, c.reader, err = c.conn.NextReader()
			if err != nil {
				return 0, err
			}
		}
		n, err := c.reader.Read(p)
		if err == io.EOF {
			// At end of message.
			c.reader = nil
			if n > 0 {
				return n, nil
			}
			// No data read, continue to next message.
			continue
		}
		return n, err
	}
}

func (c *rwc) Close() error {
	return c.conn.Close()
}

type WebSocketServer struct {
	limits   gotalk.Limits
	handlers *gotalk.Handlers

	onAccept gotalk.SockHandler

	// Template value for accepted sockets. Defaults to 0 (no automatic heartbeats)
	heartbeatInterval time.Duration

	// Template value for accepted sockets. Defaults to nil
	onHeartbeat func(load int, t time.Time)

	upgrader websocket.Upgrader
}

type WebSocketServerOptions struct {
	Limits            gotalk.Limits
	OnAccept          gotalk.SockHandler
	HeartbeatInterval time.Duration
	OnHeartbeat       func(load int, t time.Time)
	Upgrader          websocket.Upgrader
}

func DefaultWebSocketServerOptions() WebSocketServerOptions {
	return WebSocketServerOptions{
		Limits: gotalk.NewLimits(^uint32(0), ^uint32(0)),
		Upgrader: websocket.Upgrader{
			EnableCompression: true,
			HandshakeTimeout:  8 * time.Second,
			ReadBufferSize:    256,
			WriteBufferSize:   256,
			Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
				log.Errorf("websocket: status: %v, %v", status, reason)
				w.WriteHeader(status)
			},
		},
	}
}

func NewWebSocketServer(handlers *gotalk.Handlers) *WebSocketServer {
	opts := DefaultWebSocketServerOptions()
	return NewWebSocketServerWithOptions(handlers, &opts)
}

func NewWebSocketServerWithOptions(handlers *gotalk.Handlers, opts *WebSocketServerOptions) *WebSocketServer {
	ws := &WebSocketServer{
		limits:            opts.Limits,
		handlers:          handlers,
		onAccept:          opts.OnAccept,
		heartbeatInterval: opts.HeartbeatInterval,
		onHeartbeat:       opts.OnHeartbeat,
		upgrader:          opts.Upgrader,
	}
	return ws
}

func (server *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	connTakenOver := false
	defer func() {
		if !connTakenOver {
			conn.Close()
		}
	}()
	s := gotalk.NewSock(server.handlers)
	s.Adopt(&rwc{conn: conn})
	if err := s.Handshake(); err != nil {
		return err
	}
	if server.onAccept != nil {
		server.onAccept(s)
	}
	s.HeartbeatInterval = server.heartbeatInterval
	s.OnHeartbeat = server.onHeartbeat

	// Naive implementation using go routines for now.
	// TODO: Reimplement this as an event loop that handles
	// read/write for all connections with a concurrency level
	// (no. of goroutines) equal to the number of threads.

	// Start a new go-routine so that the HTTP serving stack
	// can be cleaned up. There's no need to keep the request,
	// and response writer around.
	go func() {
		defer conn.Close()
		s.Read(server.limits)
	}()
	connTakenOver = true
	return nil
}
