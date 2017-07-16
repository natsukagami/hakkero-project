package backend

// OpenSentencer represents an interface where a fetch on an opeing sentence is made.
type OpenSentencer interface {
	OpenSentence() (string, error)
}
