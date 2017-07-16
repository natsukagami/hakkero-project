package backend

import (
	"fmt"
	"log"
	"net/http"
)

// Version returns the version of Hakkero Project.
const Version = "v1.0"

// Server represents a Server.
type Server struct {
	http.Server
	c  Config
	op OpenSentencer
	q  *Queue
	r  RoomManager
}

// Welcome returns a welcome message.
func (s *Server) Welcome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(
		fmt.Sprintf("\"Welcome to Hakkero Project %s! There are %d users online waiting for a room...\"", Version, len(s.q.Players)),
	))
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}

func logRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.EscapedPath())
		h.ServeHTTP(w, r)
	})
}

func handlerApply(h http.Handler, funcs ...func(http.Handler) http.Handler) http.Handler {
	for i := len(funcs) - 1; i >= 0; i-- {
		h = funcs[i](h)
	}
	return h
}

// NewServer returns a new server.
func NewServer(c Config, op OpenSentencer, r RoomManager) *Server {
	srv := &Server{
		c:  c,
		op: op,
		r:  r,
		q:  NewQueue(r, c, op),
	}
	mux := http.NewServeMux()
	srv.Handler = handlerApply(mux, logRequest, enableCORS)
	mux.HandleFunc("/", srv.Welcome)
	mux.Handle("/rooms/", srv.r)
	mux.Handle("/queue", srv.q)
	return srv
}
