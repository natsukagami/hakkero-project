package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// The following message types are for communication by channels.

// messageIndex announces their own position on the board.
type messageIndex struct {
	Index int `json:"index"`
}

func (messageIndex) IsMessage() {}

// messageTurn contains the current turn's status, and the next turn's index and clock time.
type messageTurn struct {
	Status []Status  `json:"status"`
	Time   time.Time `json:"current"`
}

func (messageTurn) IsMessage() {}

// messageSentence contains the information of an occurring next sentence.
type messageSentence struct {
	Sentence `json:"sentence"`
	Pos      int `json:"pos"` // The position in the slice. Maybe it could be sent not in order?
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
	Content  string `json:"content"`
}

// Messager is an internal interface, to help type safety with general message-typing.
type Messager interface {
	IsMessage()
}

// Message is the general message type we will use in sending-channels.
type Message struct {
	Type    string   `json:"type"`
	Message Messager `json:"message"`
	done    chan struct{}
}

type pconnMap struct {
	mu     sync.Mutex
	Conns  map[string]*PlayerConn
	Guests []*PlayerConn
}

func (p *pconnMap) Get(id string) (*PlayerConn, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c, ok := p.Conns[id]
	return c, ok
}

func (p *pconnMap) Set(id string, conn *PlayerConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Conns[id] = conn
}

func (p *pconnMap) Guest(conn *PlayerConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Guests = append(p.Guests, conn)
}

func (p *pconnMap) Send(m Message) {
	sender := func(conn *PlayerConn) { conn.SendMessage(m) }
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, conn := range p.Conns {
		go sender(conn)
	}
	for _, conn := range p.Guests {
		go sender(conn)
	}
}

// RoomHandler is a Handler that serves players' connections to Room server.
type RoomHandler struct {
	Room Room
	// internal variables
	p         pconnMap
	ctx       context.Context
	TurnTimer *time.Timer // The turn timer.
}

// Broadcast sends the message to all listening PlayerConns.
// It pauses until all messages are scheduled to send.
func (h *RoomHandler) Broadcast(m Message) {
	h.p.Send(m)
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

func (h *RoomHandler) announceTurn(turn int) {
	h.Room.Status[turn] = StatusTurn
	sendStatus := make([]Status, len(h.Room.Status))
	copy(sendStatus, h.Room.Status)
	h.Broadcast(Message{
		Type: "turn",
		Message: messageTurn{
			Status: sendStatus,
			Time:   h.Room.Current,
		},
	})
}

// nextTurn announces the next turn, and, if ended, announces the end.
func (h *RoomHandler) nextTurn(last int) (int, bool) {
	nxt, ended := h.Room.NextTurn(last)
	if ended {
		h.Broadcast(Message{
			Type: "turn",
			Message: messageTurn{
				Status: h.Room.Status,
				Time:   h.Room.Current,
			},
		})
		h.Broadcast(Message{
			Type:    "end",
			Message: messageEnd{Winner: nxt},
		})
		return nxt, true
	}
	h.Room.Current = time.Now()
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
	h.p = pconnMap{
		Conns:  make(map[string]*PlayerConn),
		Guests: make([]*PlayerConn, 0),
	}
	var cancel context.CancelFunc
	h.ctx, cancel = context.WithCancel(context.Background())
	// Set the room up.
	h.Room = Room{
		ID:        roomID,
		Members:   players,
		Status:    make([]Status, len(players)),
		Sentences: []Sentence{Sentence{System: true, Content: openSentence}},
		Start:     time.Now(),
		Timeout:   timeout,
	}
	go h.Play(cancel)
	return
}

// Play starts up the game.
func (h *RoomHandler) Play(cancel context.CancelFunc) {
	// Wait a while so that all players are connected.
	<-time.After(10 * time.Second)
	var (
		turn  = 0
		ended = h.Room.Ended()
	)
	h.Room.Current = time.Now()
	for !ended {
		// Resets the timer so that it gives the proper time.
		h.TurnTimer = time.NewTimer(h.Room.Current.Add(h.Room.Timeout).Sub(h.Room.Current))
		conn, active := h.p.Get(h.Room.Members[turn].ID)
		h.announceTurn(turn)
		log.Printf("Room %d: Turn %d\n", h.Room.ID, turn)
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
		}
		h.TurnTimer.Stop()
		turn, ended = h.nextTurn(turn)
	}
	log.Printf("Room %d ended\n", h.Room.ID)
	cancel()
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	pConn := Prepare(conn)
	var ended bool
	select {
	case <-h.ctx.Done():
		ended = true
	default:
	}
	// If ended, immediately quit to save memory.
	if ended {
		pConn.SendMessage(Message{
			Type: "end",
			Message: messageEnd{
				Winner: h.Room.Winner(),
			},
		})
		return
	}
	// If this is a player, announce his index.
	r.ParseForm()
	index, err := h.Room.Index(r.FormValue("player"))
	if err == nil {
		ID := r.FormValue("player")
		// Replace old player connection.
		oldConn, ok := h.p.Get(ID)
		if ok {
			oldConn.Close()
		}
		pConn.SendMessage(Message{
			Type: "index",
			Message: messageIndex{
				Index: index,
			},
		})
		h.p.Set(ID, pConn)
	} else {
		// Guest,
		h.p.Guest(pConn)
	}
}
