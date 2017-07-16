package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// The following message types are for communication by channels.

// messageTurn contains the current turn's status, and the next turn's index and clock time.
type messageTurn struct {
	Status Status    `json:"status"`
	Time   time.Time `json:"time"`
	Next   int       `json:"next"`
}

func (messageTurn) IsMessage() {}

// messageSentence contains the information of an occurring next sentence.
type messageSentence struct {
	Sentence
	Pos int `json:"pos"` // The position in the slice. Maybe it could be sent not in order?
}

func (messageSentence) IsMessage() {}

// messageEnd passes the End indicator.
type messageEnd struct {
	Winner int // Index of the winner.
}

func (messageEnd) IsMessage() {}

// MessageRequest is a player's response.
type MessageRequest struct {
	IsSkip   bool `json:"skip"` // Whether the player has skipped.
	Received time.Time
	Content  string
}

// Messager is an internal interface, to help type safety with general message-typing.
type Messager interface {
	IsMessage()
}

// Message is the general message type we will use in sending-channels.
type Message struct {
	Type    string   `json:"type"`
	Message Messager `json:"message"`
}

// RoomHandler is a Handler that serves players' connections to Room server.
type RoomHandler struct {
	Room Room
	// internal variables
	Conns     map[string]*PlayerConn // All player connections.
	TurnTimer *time.Timer            // The turn timer.
}

// Broadcast sends the message to all listening PlayerConns.
// It pauses until all messages are scheduled to send.
func (h *RoomHandler) Broadcast(m Message) {
	for _, conn := range h.Conns {
		conn.Recv <- m
	}
}

// AddSentence adds a valid sentence into the Room.
func (h *RoomHandler) addSentence(id int, Content string) {
	h.Room.Status[id] = StatusActive
	sent := Sentence{
		Owner:   id,
		Content: Content,
	}
	h.Room.Sentences = append(h.Room.Sentences, sent)
	h.Broadcast(Message{
		Type: "sentence",
		Message: messageSentence{
			Sentence: sent,
			Pos:      len(h.Room.Sentences) - 1,
		},
	})
}

// addSkip adds a system skip announcement into the Room.
func (h *RoomHandler) addSkip(id int, isSkip bool) {
	sent := Sentence{System: true}
	if !isSkip {
		h.Room.Status[id] = StatusDc
		sent.Content = fmt.Sprintf("Player `%s` has timed out.", h.Room.Members[id].Username)
	} else {
		h.Room.Status[id] = StatusOut
		sent.Content = fmt.Sprintf("Player `%s` has skipped.", h.Room.Members[id].Username)
	}
	h.Room.Sentences = append(h.Room.Sentences, sent)
	h.Broadcast(Message{
		Type: "sentence",
		Message: messageSentence{
			Sentence: sent,
			Pos:      len(h.Room.Sentences) - 1,
		},
	})
}

// nextTurn announces the next turn, and, if ended, announces the end.
func (h *RoomHandler) nextTurn(last int) (int, bool) {
	nxt, ended := h.Room.NextTurn(last)
	if ended {
		h.Broadcast(Message{
			Type:    "end",
			Message: messageEnd{Winner: last},
		})
		return last, true
	}
	h.Room.Current = time.Now()
	h.Broadcast(Message{
		Type: "turn",
		Message: messageTurn{
			Status: h.Room.Status[last],
			Time:   h.Room.Current,
			Next:   nxt,
		},
	})
	return nxt, ended
}

func shufflePlayers(players []User) {
	for i := 0; i < len(players); i++ {
		nx := rand.Intn(len(players)-i) + i
		players[i], players[nx] = players[nx], players[i]
	}
}

// NewRoom creates a new room.
func NewRoom(roomID int, players []User, timeout time.Duration, openSentence string) (h *RoomHandler) {
	h = new(RoomHandler)
	shufflePlayers(players)
	h.Conns = make(map[string]*PlayerConn)
	// Set the room up.
	h.Room = Room{
		ID:        roomID,
		Members:   players,
		Status:    make([]Status, len(players)),
		Sentences: []Sentence{Sentence{System: true, Content: openSentence}},
		Start:     time.Now(),
		Current:   time.Now(),
		Timeout:   timeout,
	}
	go h.Play()
	return
}

// Play starts up the game.
func (h *RoomHandler) Play() {
	// Wait a while so that all players are connected.
	<-time.After(10 * time.Second)
	var (
		turn  = 0
		ended = h.Room.Ended()
	)
	for !ended {
		// Resets the timer so that it gives the proper time.
		h.TurnTimer = time.NewTimer(h.Room.Current.Add(h.Room.Timeout).Sub(h.Room.Current))
		conn, active := h.Conns[h.Room.Members[turn].ID]
		if !active {
			// User not even connected
			h.addSkip(turn, false)
		} else {
		awaitResp:
			for {
				select {
				case resp := <-conn.Send:
					if resp.Received.Sub(h.Room.Current) < 0 {
						continue awaitResp
					}
					if resp.IsSkip {
						h.addSkip(turn, true)
					} else if len(resp.Content) > 0 {
						h.addSentence(turn, resp.Content)
					} else {
						continue awaitResp
					}
					break awaitResp
				case err := <-conn.ErrChan:
					log.Printf("Room %d, Player %d: %v\n", h.Room.ID, turn, err)
					h.addSkip(turn, false)
					break awaitResp
				case <-h.TurnTimer.C:
					h.addSkip(turn, false)
					break awaitResp
				}
			}
			h.TurnTimer.Stop()
			turn, ended = h.nextTurn(turn)
		}
	}
	// Closes all references
	for _, conn := range h.Conns {
		conn.Close()
	}
}

func (h *RoomHandler) serveInfoReqs(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(h.Room)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\"error\": \"Server error\"}"))
		return
	}
	w.Write(data)
}

func (h *RoomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		h.serveInfoReqs(w, r)
		return
	}
	r.ParseForm()
	if _, err := h.Room.Index(r.FormValue("player")); err != nil {
		w.Write([]byte("{\"error\": \"You are not a valid player\"}"))
		w.WriteHeader(403)
		return
	}
	ID := r.FormValue("player")
	if h.Room.Ended() {
		w.Write([]byte("{\"error\": \"Game ended\"}"))
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
	pConn := Prepare(conn)
	oldConn, ok := h.Conns[ID]
	if ok {
		oldConn.Close()
	}
	h.Conns[ID] = pConn
}
