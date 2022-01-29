package websocks

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/akutz/memconn"
	socks5 "github.com/armon/go-socks5"
	"golang.org/x/net/websocket"
)

// Websocketに流れてきたパケットをsocks5として解釈するサーバー
type Server struct {
	URI   string
	spath string
	host  string

	http.ServeMux
	Conf socks5.Config
}

func (s *Server) parseHost() (err error) {
	if !strings.Contains(s.URI, "//") {
		s.URI = "ws://" + s.URI
	}

	loc, err := url.ParseRequestURI(s.URI)
	if err != nil {
		return
	}
	s.spath = loc.Path
	s.host = loc.Host

	if s.spath == "" {
		s.spath = "/"
	}

	return nil
}

func (s *Server) ListenAndServe() (err error) {
	if err = s.parseHost(); err != nil {
		return
	}

	if s.host == "" {
		s.host = ":80"
	}
	ln, err := net.Listen("tcp", s.host)
	if err != nil {
		return
	}
	return s.Serve(ln)
}

func (s *Server) Serve(li net.Listener) (err error) {
	if err = s.parseHost(); err != nil {
		return
	}

	soc, err := socks5.New(&s.Conf)
	if err != nil {
		return
	}

	lis, err := memconn.Listen("memu", "sockListener")
	if err != nil {
		return
	}

	s.Handle(s.spath, websocket.Handler(func(ws *websocket.Conn) {
		sock, err := memconn.Dial("memu", "sockListener")
		if err != nil {
			return
		}
		go func() {
			io.Copy(sock, ws)
			ws.Close()
		}()
		io.Copy(ws, sock)
		sock.Close()
	}))

	go soc.Serve(lis) // socks5 server
	if err := http.Serve(li, s); err != nil {
		return err
	}
	return nil
}

func (s *Server) Path() string {
	return s.spath
}

func (s *Server) Host() string {
	return s.host
}
