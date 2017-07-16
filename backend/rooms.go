package backend

import (
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Rooms implement a Room container.
type Rooms struct {
	mu    sync.Mutex
	Rooms []*RoomHandler // As an array, for easy numbering and in-memory management.
}

// New creates a new room(handler) and return its id.
func (r *Rooms) New(players []User, timeout time.Duration, openSentence string) (id int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Rooms == nil {
		r.Rooms = make([]*RoomHandler, 0)
	}
	r.Rooms = append(r.Rooms, NewRoom(len(r.Rooms), players, timeout, openSentence))
	return len(r.Rooms) - 1, nil
}

// Get returns the room with the specified ID.
func (r *Rooms) Get(id int) (*RoomHandler, error) {
	if len(r.Rooms) <= id {
		return nil, errors.New("no such room")
	}
	return r.Rooms[id], nil
}

func (r *Rooms) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	if len(rq.URL.EscapedPath()) <= len("/rooms/") {
		w.WriteHeader(403)
		return
	}
	idStr := rq.URL.EscapedPath()[len("/rooms/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	room, err := r.Get(id)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(err.Error()))
		return
	}
	room.ServeHTTP(w, rq)
}

// RoomManager is an interface for a RoomManager.
type RoomManager interface {
	http.Handler
	New(players []User, timeout time.Duration, openSentence string) (id int, err error)
	Get(id int) (*RoomHandler, error)
}

var _ RoomManager = (*Rooms)(nil)
