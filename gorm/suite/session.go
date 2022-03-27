package suite

// Copied from https://github.com/gobuffalo/suite/blob/master/session.go
// See https://github.com/AlphaFlow/api-core/pull/119
// See https://github.com/AlphaFlow/institutional-api/pull/1444

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionStore struct {
	sessions map[string]*sessions.Session
}

func (s *sessionStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	if s, ok := s.sessions[name]; ok {
		return s, nil
	}
	return s.New(r, name)
}

func (s *sessionStore) New(r *http.Request, name string) (*sessions.Session, error) {
	sess := sessions.NewSession(s, name)
	s.sessions[name] = sess
	return sess, nil
}

func (s *sessionStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error {
	if s.sessions == nil {
		s.sessions = map[string]*sessions.Session{}
	}
	s.sessions[sess.Name()] = sess
	return nil
}

//NewSessionStore for action suite
func newSessionStore() sessions.Store {
	return &sessionStore{
		sessions: map[string]*sessions.Session{},
	}
}
