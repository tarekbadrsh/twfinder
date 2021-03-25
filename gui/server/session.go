// Implementation of the GUI session.

package server

import (
	"crypto/rand"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
	"twfinder/logger"
)

// Session interface defines the session to the GUI users (clients).
type Session interface {
	// ID returns the ID of the session.
	ID() string

	// New tells if the session is new meaning the client
	// does not (yet) know about it.
	New() bool

	// Private tells if the session is a private session.
	// There is only one public session, and it is shared
	// between the "sessionless" users.
	Private() bool

	// AddWin adds a window to the session.
	// Returns an error if window name is empty or
	// a window with the same name has already been added.
	AddWin(w Window) error

	// RemoveWin removes a window from the session.
	// Returns if the window was removed from the session.
	RemoveWin(w Window) bool

	// SortedWins returns a sorted slice of windows.
	// The slice is sorted by window text (title).
	SortedWins() []Window

	// WinByName returns a window specified by its name.
	WinByName(name string) Window

	// Attr returns the value of an attribute stored in the session.
	// TODO use an interface type something like "serializable".
	Attr(name string) interface{}

	// SetAttr sets the value of an attribute stored in the session.
	// Pass the nil value to delete the attribute.
	SetAttr(name string, value interface{})

	// Created returns the time when the session was created.
	Created() time.Time

	// Accessed returns the time when the session was last accessed.
	Accessed() time.Time

	// Timeout returns the session timeout.
	Timeout() time.Duration

	// SetTimeout sets the session timeout.
	SetTimeout(timeout time.Duration)

	// access registers an access to the session.
	// Implementation locks or the sessions RW mutex.
	access()

	// ClearNew clears the new flag.
	// After this New() will return false.
	clearNew()

	// rwMutex returns the RW mutex of the session.
	rwMutex() *sync.RWMutex
}

// Session implementation.
type sessionImpl struct {
	id       string                 // ID of the session
	isNew    bool                   // Tells if the session is new
	created  time.Time              // Creation time
	accessed time.Time              // Last accessed time
	windows  map[string]Window      // Windows of the session
	attrs    map[string]interface{} // Attributes stored in the session
	timeout  time.Duration          // Session timeout

	rwMutexF *sync.RWMutex // RW mutex to synchronize session (and related Window and component) access
}

// newSessionImpl creates a new sessionImpl.
// The default timeout is 30 minutes.
func newSessionImpl(private bool) sessionImpl {
	var id string
	// The public session has an empty string ID
	if private {
		id = genID()
	}

	now := time.Now()

	// Initialzie private sessions as new, but not the public session
	return sessionImpl{id: id, isNew: private, created: now, accessed: now, windows: make(map[string]Window),
		attrs: make(map[string]interface{}), timeout: 30 * time.Minute, rwMutexF: &sync.RWMutex{}}
}

// Valid characters (bytes) to be used in session IDs
// Its length must be a power of 2.
const idChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"

func init() {
	// Is len(idChars) a power of 2?
	if i := byte(len(idChars)); i&(i-1) != 0 {
		panic(fmt.Sprint("len(idChars) must be power of 2: ", i))
	}
}

// Length of the session IDs
const idLength = 22

// genID generates a new session ID.
func genID() string {
	id := make([]byte, idLength)
	if _, err := rand.Read(id); err != nil {
		logger.Errorf("Failed to read from secure random: %v", err)
	}

	for i, v := range id {
		id[i] = idChars[v&byte(len(idChars)-1)]
	}
	return string(id)
}

func (s *sessionImpl) ID() string {
	return s.id
}

func (s *sessionImpl) New() bool {
	return s.isNew
}

func (s *sessionImpl) Private() bool {
	return len(s.id) > 0
}

func (s *sessionImpl) AddWin(w Window) error {
	if len(w.Name()) == 0 {
		return errors.New("Window name cannot be empty string")
	}
	if _, exists := s.windows[w.Name()]; exists {
		return errors.New("A window with the same name has already been added: " + w.Name())
	}

	s.windows[w.Name()] = w

	return nil
}

func (s *sessionImpl) RemoveWin(w Window) bool {
	win := s.windows[w.Name()]
	if win != nil && win.ID() == w.ID() {
		delete(s.windows, w.Name())
		return true
	}
	return false
}

func (s *sessionImpl) SortedWins() []Window {
	wins := make(WinSlice, len(s.windows))

	i := 0
	for _, win := range s.windows {
		wins[i] = win
		i++
	}

	sort.Sort(wins)

	return wins
}

func (s *sessionImpl) WinByName(name string) Window {
	return s.windows[name]
}

func (s *sessionImpl) Attr(name string) interface{} {
	return s.attrs[name]
}

func (s *sessionImpl) SetAttr(name string, value interface{}) {
	if value == nil {
		delete(s.attrs, name)
	} else {
		s.attrs[name] = value
	}
}

func (s *sessionImpl) Created() time.Time {
	return s.created
}

func (s *sessionImpl) Accessed() time.Time {
	s.rwMutexF.RLock()
	defer s.rwMutexF.RUnlock()
	return s.accessed
}

func (s *sessionImpl) Timeout() time.Duration {
	return s.timeout
}

func (s *sessionImpl) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

func (s *sessionImpl) access() {
	s.rwMutexF.Lock()
	s.accessed = time.Now()
	s.rwMutexF.Unlock()
}

func (s *sessionImpl) clearNew() {
	s.isNew = false
}

func (s *sessionImpl) rwMutex() *sync.RWMutex {
	return s.rwMutexF
}
