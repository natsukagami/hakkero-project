package backend

import (
	"encoding/json"
	"math/rand"
)

// User represents an active user. The user can be in a room, or in the placement queue.
// In this implementation we will NOT save user information into the database, as we allow
// register-less entrances. Therefore, we would like User to be as simple as possible.
type User struct {
	ID       string
	Username string
}

// MarshalJSON turns an user into a JSON string.
// As the user struct should not expose an user's ID, only the Username should be sent.
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Username)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// NewUser creates a new unique user.
func NewUser(username string) User {
	return User{
		Username: username,
		ID:       randStringBytesMaskImpr(32),
	}
}
