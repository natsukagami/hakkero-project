package backend

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageQueueResponse represents an user's answer to the queuing request.
type MessageQueueResponse struct {
	Accepted bool `json:"accepted"`
	Received time.Time
}

// messageQueueFound represents a "Match Found" message.
type messageQueueFound struct{}

func (m messageQueueFound) IsMessage() {}

// messageQueueAnnouncement represents a "Matchmaking Success/Failed" message.
type messageQueueAnnouncement struct {
	Success      bool   `json:"success"`
	Room         int    `json:"room,omitempty"`
	Announcement string `json:"announcement,omitempty"`
}

func (m messageQueueAnnouncement) IsMessage() {}

// messageQueueSize sends any queue size update.
type messageQueueSize struct {
	Size int `json:"size"`
}

func (m messageQueueSize) IsMessage() {}

// Queue represents a queue handler.
type Queue struct {
	Config  Config
	Rooms   RoomManager
	OP      OpenSentencer
	mu      sync.Mutex
	Players []*QueueConn
}

// Broadcast sends a message to all audiences.
func (q *Queue) Broadcast(audience []*QueueConn, m Message) {
	for _, conn := range audience {
		conn.Recv <- m
	}
}

func (q *Queue) awaitResponse(player *QueueConn, timeout time.Duration) (accept bool, received bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case m := <-player.Send:
			return m.Accepted, true
		case <-timer.C:
			return false, false
		}
	}
}

func userFromConn(arr []*QueueConn) []User {
	ans := make([]User, len(arr))
	for id, conn := range arr {
		ans[id] = conn.User
	}
	return ans
}

// Play prompts the players for a game.
// It returns the players who accepted but not enough for a game.
func (q *Queue) Play(players []*QueueConn) []*QueueConn {
	q.Broadcast(players, Message{
		Type:    "found",
		Message: messageQueueFound{},
	})
	accepted := make(chan *QueueConn)
	acceptedArr := make([]*QueueConn, 0)
	for _, player := range players {
		go func(p *QueueConn) {
			accept, received := q.awaitResponse(p, 10*time.Second)
			if accept {
				accepted <- p
				return
			}
			if !received {
				p.Recv <- Message{
					Type: "announcement",
					Message: messageQueueAnnouncement{
						Success:      false,
						Announcement: "You have timed-out a game. Please refresh the page to join again.",
					},
				}
			} else {
				p.Recv <- Message{
					Type: "announcement",
					Message: messageQueueAnnouncement{
						Success:      false,
						Announcement: "You have rejected a game. Please refresh the page to join again.",
					},
				}
			}
			p.Close()
			accepted <- nil
		}(player)
	}
	for i := 0; i < len(players); i++ {
		recv := <-accepted
		if recv != nil {
			acceptedArr = append(acceptedArr, recv)
		}
	}
	if len(acceptedArr) == len(players) {
		os, err := q.OP.OpenSentence()
		if err != nil {
			q.Broadcast(acceptedArr, Message{
				Type: "announcement",
				Message: messageQueueAnnouncement{
					Success:      false,
					Announcement: "Cannot find a proper open sentence. Match cancelled!",
				},
			})
			return acceptedArr
		}
		// New game accepted.
		id, err := q.Rooms.New(userFromConn(acceptedArr), q.Config.Timeout, os)
		if err != nil {
			q.Broadcast(acceptedArr, Message{
				Type: "announcement",
				Message: messageQueueAnnouncement{
					Success:      false,
					Announcement: "Cannot set up a game room. Match cancelled!",
				},
			})
			return acceptedArr
		}
		q.Broadcast(acceptedArr, Message{
			Type: "announcement",
			Message: messageQueueAnnouncement{
				Success:      true,
				Room:         id,
				Announcement: fmt.Sprintf("You have been assigned to room %d. Match starting soon!", id),
			},
		})
		return nil
	}
	q.Broadcast(acceptedArr, Message{
		Type: "announcement",
		Message: messageQueueAnnouncement{
			Success:      false,
			Announcement: "Match cannot start because someone failed the ready check.",
		},
	})
	return acceptedArr
}

// Enqueue adds a player into the queue.
func (q *Queue) Enqueue(player *QueueConn) {
	q.mu.Lock()
	q.Players = append(q.Players, player)
	q.Broadcast(q.Players, Message{
		Type: "size",
		Message: messageQueueSize{
			Size: len(q.Players),
		},
	})
	q.mu.Unlock()
	if len(q.Players) == q.Config.PlayerLimit {
		old := q.Players
		q.Players = make([]*QueueConn, 0)
		for _, player := range q.Play(old) {
			q.Enqueue(player)
		}
	}
}

func (q *Queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	if len(username) == 0 || len(username) > 20 {
		w.Write([]byte("{\"error\": \"Invalid username\"}"))
		w.WriteHeader(400)
		return
	}
	upg := websocket.Upgrader{}
	conn, err := upg.Upgrade(w, r, nil)
	if err != nil {
		w.Write([]byte("{\"error\": \"Cannot upgrade to WebSocket\"}"))
		w.WriteHeader(400)
		return
	}
	pConn := Enqueue(conn, username)
	q.Enqueue(pConn)
}

// NewQueue returns a new queue.
func NewQueue(r RoomManager, c Config, op OpenSentencer) *Queue {
	return &Queue{
		Config:  c,
		Rooms:   r,
		OP:      op,
		Players: make([]*QueueConn, 0),
	}
}
