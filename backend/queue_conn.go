package backend

import (
	"time"

	"github.com/gorilla/websocket"

	"github.com/pkg/errors"
)

// QueueConn represents a connection to the queue.
type QueueConn struct {
	Conn
	User
	Send chan MessageQueueResponse
}

// forwarder fetches messages from user interface and forwards it to Handler.
func (q *QueueConn) forwarder() {
	q.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	defer close(q.Send)
	for q.Error == nil {
		var ms MessageQueueResponse
		err := q.Conn.ReadJSON(&ms)
		if err != nil {
			go q.broadcastError(errors.Wrap(err, "playerconn read"))
			return
		}
		ms.Received = time.Now()
		q.Send <- ms
	}
}

// Enqueue fires up the QueueConn for usage
func Enqueue(conn *websocket.Conn, username string) *QueueConn {
	q := &QueueConn{
		Conn: Conn{
			Conn:    conn,
			ErrChan: make(chan error)},
		User: NewUser(username),
	}
	q.Recv = make(chan Message)
	q.Send = make(chan MessageQueueResponse)
	go q.forwarder()
	go q.receiver()
	return q
}
