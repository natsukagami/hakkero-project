package backend

// Sentence is a sentence written is a room-wide paragraph.
// The sentence could be an user's sentence, or a system announcement
// (e.g. An user has left the game).
type Sentence struct {
	Content string `json:"content"`
	Owner   int    `json:"owner"`           // The user index (in the User slice) who wrote this sentence. In the case of a system announcement, this is left empty.
	System  bool   `json:"system,omitempy"` // Indicate that it's a system announcement.
}
