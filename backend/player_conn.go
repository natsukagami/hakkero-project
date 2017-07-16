package backend

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Conn represents a client connection.
type Conn struct {
	*websocket.Conn
	Recv chan Message // The channel for receiving messages.

	Error error // The error variable, if it is set then the connection no longer is valuable.
	// the dedicated error channel
	ErrChan chan error
}

// PlayerConn represents a Player Connection.
type PlayerConn struct {
	Conn
	Send chan MessageRequest // The channel for sending messages.
}

// Just broadcast error forever.
func (p *Conn) broadcastError(err error) {
	if p.Error != nil {
		return // Only do the first error
	}
	p.Error = err
	for {
		p.ErrChan <- p.Error
	}
}

// receiver fetches messages from Handler and passes it to user.
func (p *Conn) receiver() {
	for ms := range p.Recv {
		if p.Error != nil {
			continue // If an error occurs, just consume.
		}
		err := p.Conn.WriteJSON(&ms)
		if err != nil {
			go p.broadcastError(errors.Wrap(err, "playerconn write"))
		}
	}
}

// forwarder fetches messages from user interface and forwards it to Handler.
func (p *PlayerConn) forwarder() {
	p.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	defer close(p.Send)
	for p.Error == nil {
		var ms MessageRequest
		err := p.Conn.ReadJSON(&ms)
		if err != nil {
			go p.broadcastError(errors.Wrap(err, "playerconn read"))
			return
		}
		ms.Received = time.Now()
		p.Send <- ms
	}
}

// Close closes the player connection.
func (p *Conn) Close() error {
	close(p.Recv)
	return p.Conn.Close()
}

// Prepare fires up the PlayerConn for usage
func Prepare(conn *websocket.Conn) *PlayerConn {
	p := &PlayerConn{
		Conn: Conn{
			Conn:    conn,
			ErrChan: make(chan error)},
	}
	p.Recv = make(chan Message)
	p.Send = make(chan MessageRequest)
	go p.forwarder()
	go p.receiver()
	return p
}
