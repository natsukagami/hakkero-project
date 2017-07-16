package backend

import (
	"encoding/json"
	"errors"
	"time"
)

// Status represents a player's status in a room.
type Status int

// Status constants.
const (
	StatusActive Status = iota
	StatusOut
	StatusDc
	StatusTurn
)

// MarshalJSON returns the status of the player in JSON.
// Numbers are not the pretty thing, we send status strings instead.
func (s Status) MarshalJSON() ([]byte, error) {
	var st string
	switch s {
	case StatusActive:
		st = "active"
	case StatusOut:
		st = "skipped"
	case StatusDc:
		st = "disconnected"
	case StatusTurn:
		st = "turn"
	default:
		st = "unknown"
	}
	return json.Marshal(st)
}

// Room represents a playing Room.
type Room struct {
	ID        int           `json:"id"`
	Members   []User        `json:"members"`   // The list of all members.
	Status    []Status      `json:"status"`    // The status of each member in the room.
	Sentences []Sentence    `json:"sentences"` // The list of all sentences.
	Start     time.Time     `json:"start"`     // The time of the start of the game.
	Current   time.Time     `json:"current"`   // The time when the current turn started.
	Timeout   time.Duration `json:"timeout"`   // The time of each turn.
}

// Ended returns whether the game has ended.
func (r Room) Ended() bool {
	// Simply put, the game ended if and only if only one player is active.
	active := 0
	for _, status := range r.Status {
		if status == StatusActive || status == StatusTurn {
			active++
		}
	}
	return active == 1
}

// Index returns a player's index in the slice.
// Throws an error if it's not found.
func (r Room) Index(ID string) (int, error) {
	for id, user := range r.Members {
		if user.ID == ID {
			return id, nil
		}
	}
	return 0, errors.New("player not found")
}

// NextTurn returns the next turn and DOES NOT modifies the current status.
// Returns false if game ended.
func (r Room) NextTurn(last int) (int, bool) {
	if r.Ended() {
		return last, false
	}
	cur := last
	for {
		cur = (cur + 1) % len(r.Members)
		if r.Status[cur] == StatusActive {
			return cur, true
		}
	}
}
