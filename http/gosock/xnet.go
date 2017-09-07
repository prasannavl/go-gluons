package gosock

import (
	"errors"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/prasannavl/go-gluons/log"

	"github.com/rsms/gotalk"
)

type XNetWebSocketServer struct {
	limits   gotalk.Limits
	handlers *gotalk.Handlers

	onAccept gotalk.SockHandler

	// Template value for accepted sockets. Defaults to 0 (no automatic heartbeats)
	heartbeatInterval time.Duration

	// Template value for accepted sockets. Defaults to nil
	onHeartbeat func(load int, t time.Time)

	Server *websocket.Server
}

func NewXNetWebSocketServer(handlers *gotalk.Handlers, limits gotalk.Limits, onAccept gotalk.SockHandler) *XNetWebSocketServer {
	ws := &XNetWebSocketServer{
		limits:   limits,
		handlers: handlers,
		onAccept: onAccept,
	}
	ws.Server = &websocket.Server{Handler: ws.handleConn, Handshake: checkOrigin}
	return ws
}

func (server *XNetWebSocketServer) handleConn(ws *websocket.Conn) {
	s := gotalk.NewSock(server.handlers)
	ws.PayloadType = websocket.BinaryFrame // websocket.TextFrame;
	s.Adopt(ws)
	if err := s.Handshake(); err != nil {
		log.Tracef("gosock: %v", err)
		return
	}
	if server.onAccept != nil {
		server.onAccept(s)
	}
	s.HeartbeatInterval = server.heartbeatInterval
	s.OnHeartbeat = server.onHeartbeat
	s.Read(server.limits)
}

func (s *XNetWebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	// This hijacks the conn, defers a close,
	// and call the websocket handler which is the above handler
	s.Server.ServeHTTP(w, r)
	return nil
}

func checkOrigin(config *websocket.Config, req *http.Request) (err error) {
	config.Origin, err = websocket.Origin(config, req)
	if err == nil && config.Origin == nil {
		return errors.New("gosock: null origin")
	}
	return err
}
