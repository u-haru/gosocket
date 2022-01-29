package websocks

import (
	"errors"
	"io"
	"log"
	"net"
	"net/url"

	"golang.org/x/net/websocket"
)

// HostへのパケットをwebsocketでTargetに流すクライアント
type Client struct {
	Target   string
	Host     string
	url      string
	protocol string
	origin   string
}

func (c *Client) ListenAndServe() (err error) {
	if c.Host == "" {
		c.Host = ":80"
	}
	if c.Target == "" {
		return errors.New("target isn't specified")
	}
	ln, err := net.Listen("tcp", c.Host)
	if err != nil {
		return
	}
	return c.Serve(ln)
}

func (c *Client) Serve(li net.Listener) (err error) {
	loc, err := url.ParseRequestURI(c.Target)
	if err != nil {
		return
	}
	c.url = loc.Scheme + "://" + loc.Host + loc.Path
	c.protocol = loc.Scheme
	c.origin = loc.Scheme + "://" + loc.Host

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		ws, err := websocket.Dial(c.url, c.protocol, c.origin)
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		go func() {
			io.Copy(conn, ws)
			ws.Close()
			conn.Close() // wsが閉じてもconnは勝手に閉じてくれない
		}()
		go func() {
			io.Copy(ws, conn)
			conn.Close()
			ws.Close() // 一応閉じとく
		}()
	}
}

func (c *Client) Url() string {
	return c.url
}

func (c *Client) Protocol() string {
	return c.protocol
}
func (c *Client) Origin() string {
	return c.origin
}
