// Implementation of the GUI server which handles sessions,
// renders the windows and handles event dispatching.

package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Internal path constants.
const (
	pathStatic     = "_gui_static/" // App path-relative path for GUI static contents.
	pathSessCheck  = "_sess_ch"     // App path-relative path for checking session (without registering access)
	pathEvent      = "e"            // Window-relative path for sending events
	pathRenderComp = "rc"           // Window-relative path for rendering a component
)

// Parameters passed between the browser and the server.
const (
	paramEventType     = "et"   // Event type parameter name
	paramCompID        = "cid"  // Component id parameter name
	paramCompValue     = "cval" // Component value parameter name
	paramFocusedCompID = "fcid" // Focused component id parameter name
	paramMouseWX       = "mwx"  // Mouse x pixel coordinate (inside window)
	paramMouseWY       = "mwy"  // Mouse y pixel coordinate (inside window)
	paramMouseX        = "mx"   // Mouse x pixel coordinate (relative to source component)
	paramMouseY        = "my"   // Mouse y pixel coordinate (relative to source component)
	paramMouseBtn      = "mb"   // Mouse button
	paramModKeys       = "mk"   // Modifier key states
	paramKeyCode       = "kc"   // Key code
)

// Event response actions (client actions to take after processing an event).
const (
	eraNoAction   = iota // Event processing OK and no action required
	eraReloadWin         // Window name to be reloaded
	eraDirtyComps        // There are dirty components which needs to be refreshed
	eraFocusComp         // Focus a component
)

// Default GUI session id cookie name
const defaultSessIDCookieName = "gui-sessid"

// SessionHandler interface defines a callback to get notified
// for certain events related to session life-cycles.
type SessionHandler interface {
	// Created is called when a new session is created.
	// At this time the client does not yet know about the session.
	Created(sess Session)

	// Removed is called when a session is being removed
	// from the server. After removal, the session id will become
	// an invalid session id.
	Removed(sess Session)
}

// AppRootHandlerFunc is the function type that handles the application root (when no window name is specified).
// sess is the shared, public session if no private session is created.
type AppRootHandlerFunc func(w http.ResponseWriter, r *http.Request, sess Session)

// Server interface defines the GUI server which handles sessions,
// renders the windows, components and handles event dispatching.
type Server interface {
	// The Server implements the Session interface:
	// there is one public session which is shared between
	// the "sessionless" requests.
	// This is to maintain windows without a session.
	Session

	// A server has text which will be used as the title
	// of the server.
	HasText

	// Secure returns if the server is configured to run
	// in secure (HTTPS) mode or in HTTP mode.
	Secure() bool

	// AppURL returns the application URL string.
	AppURL() string

	// AppPath returns the application path string.
	AppPath() string

	// AddSessCreatorName registers a nonexistent window name
	// whose path auto-creates a new session.
	//
	// Normally sessions are created from event handlers during
	// event dispatching by calling Event.NewSession(). This
	// requires a public window and an event source component
	// (e.g. a Button) to create a session.
	// With AddSessCreatorName you can register nonexistent (meaning
	// not-yet added) window names whose path will trigger an automatic
	// session creation (if the current session is not private), and
	// with a registered SessionHandler you can build the window and
	// add it to the auto-created new session prior to it being served.
	//
	// The text linking to the name will be included in the window list
	// if text is a non-empty string.
	//
	// Tip: You can use this to pre-register a login window for example.
	// You can call
	// 		AddSessCreatorName("login", "Login Window")
	// and in the Created() method of a registered SessionHandler:
	//		func (h MySessHanlder) Created(s server.Session) {
	//			win := server.NewWindow("login", "Login Window")
	//			// ...add content to the login window...
	// 			s.AddWindow(win)
	// 		}
	AddSessCreatorName(name, text string)

	// AddSHandler adds a new session handler.
	AddSHandler(handler SessionHandler)

	// SetHeaders sets extra HTTP response headers that are added to all responses.
	// Supplied values are copied, so changes to the passed map afterwards have no effect.
	//
	// For example to add an extra "Server" header whose value is the server version:
	//     server.SetHeaders(map[string][]string{
	//         "Server": {server.ServerVersion},
	//     })
	SetHeaders(headers map[string][]string)

	// Headers returns the extra HTTP response headers that are added to all repsonses.
	// A copy is returned, so changes to the returned map afterwards have no effect.
	Headers() map[string][]string

	// AddStaticDir registers a directory whose content (files) recursively
	// will be served by the server when requested.
	// path is an app-path relative path to address a file, dir is the root directory
	// to search in.
	// Note that the app name must be included in absolute request paths,
	// and it may be omitted if you want to use relative paths.
	// Extra headers set by SetHeaders() will also be included in responses serving the static files.
	//
	// Example:
	//     AddStaticDir("img", "/tmp/myimg")
	// Then request for absolute path "/appname/img/faces/happy.gif" will serve
	// "/tmp/myimg/faces/happy.gif", just as the the request for relative path "img/faces/happy.gif".
	AddStaticDir(path, dir string) error

	// Theme returns the default CSS theme of the server.
	Theme() string

	// SetTheme sets the default CSS theme of the server.
	SetTheme(theme string)

	// SetLogger sets the logger to be used
	// to log incoming requests.
	// Pass nil to disable logging. This is the default.
	SetLogger(logger *log.Logger)

	// Logger returns the logger that is used to log incoming requests.
	Logger() *log.Logger

	// AddRootHeadHTML adds an HTML text which will be included
	// in the HTML <head> section of the window list page (the app root).
	// Note that these will be ignored if you take over the app root
	// (by calling SetAppRootHandler).
	AddRootHeadHTML(html string)

	// RemoveRootHeadHTML removes an HTML head text
	// that was previously added with AddRootHeadHTML().
	RemoveRootHeadHTML(html string)

	// SetAppRootHandler sets a function that is called when the app root is requested.
	// The default function renders the window list, including authenticated windows
	// and session creators - with clickable links.
	// By setting your own hander, you will completely take over the app root.
	SetAppRootHandler(f AppRootHandlerFunc)

	// SessIDCookieName returns the cookie name used to store the sever
	// session ID.
	SessIDCookieName() string

	// session ID.
	SetSessIDCookieName(name string)

	// Start starts the GUI server and waits for incoming connections.
	//
	// Sessionless window names may be specified as optional parameters
	// that will be opened in the default browser.
	// Tip: Pass an empty string to open the window list.
	// Tip: Not passing any window names will start the server silently
	// without opening any windows.
	Start(openWins ...string) error
}

// Server implementation.
type serverImpl struct {
	sessionImpl // Single public session implementation
	hasTextImpl // Has text implementation

	appName            string             // Application name (part of the application path)
	addr               string             // Server address
	secure             bool               // Tells if the server is configured to run in secure (HTTPS) mode
	appPath            string             // Application path
	appURLString       string             // Application URL string
	appURL             *url.URL           // Application URL, parsed
	sessions           map[string]Session // Sessions
	certFile, keyFile  string             // Certificate and key files for secure (HTTPS) mode
	sessCreatorNames   map[string]string  // Session creator names
	sessionHandlers    []SessionHandler   // Registered session handlers
	theme              string             // Default CSS theme of the server
	logger             *log.Logger        // Logger.
	headers            http.Header        // Extra headers that will be added to all responses.
	rootHeads          []string           // Additional head HTML texts of the window list page (app root)
	appRootHandlerFunc AppRootHandlerFunc // App root handler function
	sessIDCookieName   string             // Session ID cookie name

	sessMux sync.RWMutex // Mutex to protect state related to session handling
}

// NewServer creates a new GUI server in HTTP mode.
// The specified app name will be part of the application path (the first part).
// If addr is empty string, "localhost:3434" will be used.
//
// Tip: Pass an empty string as appName to place the GUI server to the root path ("/").
func NewServer(appName, addr string) Server {
	return newServerImpl(appName, addr, "", "")
}

// NewServerTLS creates a new GUI server in secure (HTTPS) mode.
// The specified app name will be part of the application path (the first part).
// If addr is empty string, "localhost:3434" will be used.
//
// Tip: Pass an empty string as appName to place the GUI server to the root path ("/").
// Tip: You can use generate_cert.go in crypto/tls to generate
// a test certificate and key file (cert.pem andkey.pem).
func NewServerTLS(appName, addr, certFile, keyFile string) Server {
	return newServerImpl(appName, addr, certFile, keyFile)
}

// newServerImpl creates a new serverImpl.
func newServerImpl(appName, addr, certFile, keyFile string) *serverImpl {
	if addr == "" {
		addr = "localhost:3434"
	}

	s := &serverImpl{
		sessionImpl:      newSessionImpl(false),
		appName:          appName,
		addr:             addr,
		sessions:         make(map[string]Session),
		sessCreatorNames: make(map[string]string),
		theme:            ThemeDefault,
		sessIDCookieName: defaultSessIDCookieName,
	}

	if s.appName == "" {
		s.appPath = "/"
	} else {
		s.appPath = "/" + s.appName + "/"
	}

	if certFile == "" || keyFile == "" {
		s.secure = false
		s.appURLString = "http://" + addr + s.appPath
	} else {
		s.secure = true
		s.appURLString = "https://" + addr + s.appPath
		s.certFile = certFile
		s.keyFile = keyFile
	}
	var err error
	if s.appURL, err = url.Parse(s.appURLString); err != nil {
		panic(fmt.Sprintf("Parse %q: %+v", s.appURLString, err))
	}

	s.appRootHandlerFunc = s.renderWinList

	return s
}

func (s *serverImpl) Secure() bool {
	return s.secure
}

func (s *serverImpl) AppURL() string {
	return s.appURLString
}

func (s *serverImpl) AppPath() string {
	return s.appPath
}

func (s *serverImpl) AddSessCreatorName(name, text string) {
	if len(name) > 0 {
		s.sessCreatorNames[name] = text
	}
}

func (s *serverImpl) AddSHandler(handler SessionHandler) {
	s.sessMux.Lock()
	s.sessionHandlers = append(s.sessionHandlers, handler)
	s.sessMux.Unlock()
}

// newSession creates a new (private) Session.
// The event is optional. If specified and the current session
// (as returned by Event.Session()) is private, it will be removed first.
// The new session is set to the event, and also returned.
func (s *serverImpl) newSession(e *eventImpl) Session {
	if e != nil {
		// First remove old session
		s.removeSess(e)
	}

	sessImpl := newSessionImpl(true)
	sess := &sessImpl
	if e != nil {
		e.shared.session = sess
	}
	// Store new session
	s.sessMux.Lock()
	s.sessions[sess.ID()] = sess

	if s.logger != nil {
		s.logger.Println("SESSION created:", sess.ID())
	} else {
		log.Println("SESSION created:", sess.ID())
	}

	// Notify session handlers
	for _, handler := range s.sessionHandlers {
		handler.Created(sess)
	}
	s.sessMux.Unlock()

	return sess
}

// removeSess removes (invalidates) the current session of the specified event.
// Only private sessions can be removed, calling this
// when the current session (as returned by Event.Session()) is public is a no-op.
// After this method Event.Session() will return the shared public session.
func (s *serverImpl) removeSess(e *eventImpl) {
	if e.shared.session.Private() {
		s.sessMux.Lock()
		s.removeSess2(e.shared.session)
		s.sessMux.Unlock()
		e.shared.session = &s.sessionImpl
	}
}

// removeSess2 removes (invalidates) the specified session.
// Only private sessions can be removed, calling this with the
// public session is a no-op.
// serverImpl.mux must be locked when this is called.
func (s *serverImpl) removeSess2(sess Session) {
	if sess.Private() {
		if s.logger != nil {
			s.logger.Println("SESSION removed:", sess.ID())
		} else {
			log.Println("SESSION removed:", sess.ID())
		}

		// Notify session handlers
		for _, handler := range s.sessionHandlers {
			handler.Removed(sess)
		}
		delete(s.sessions, sess.ID())
	}
}

// addSessCookie lets the client know about the specified (new) session
// by setting the GUI session id cookie.
// Also clears the new flag of the session.
func (s *serverImpl) addSessCookie(sess Session, w http.ResponseWriter) {
	// HttpOnly: do not allow non-HTTP access to it (like javascript) to prevent stealing it...
	// Secure: only send it over HTTPS
	// MaxAge: to specify the max age of the cookie in seconds, else it's a session cookie and gets deleted after the browser is closed.
	c := http.Cookie{
		Name:     s.sessIDCookieName,
		Value:    sess.ID(),
		Path:     s.appURL.EscapedPath(),
		HttpOnly: true,
		Secure:   s.secure,
		MaxAge:   72 * 60 * 60, // 72 hours max age
	}
	http.SetCookie(w, &c)

	sess.clearNew()
}

// sessCleaner periodically checks whether private sessions has timed out
// in an endless loop. If a session has timed out, removes it.
// This method is to start as a new go routine.
func (s *serverImpl) sessCleaner() {
	sleep := 10 * time.Second
	for {
		now := time.Now()

		s.sessMux.Lock()
		for _, sess := range s.sessions {
			if now.Sub(sess.Accessed()) > sess.Timeout() {
				s.removeSess2(sess)
			}
		}
		s.sessMux.Unlock()

		time.Sleep(sleep)
	}
}

func (s *serverImpl) SetHeaders(headers map[string][]string) {
	s.headers = make(map[string][]string, len(headers))
	for k, v := range headers {
		// Also copy value which is a slice
		s.headers[k] = append(make([]string, 0, len(v)), v...)
	}
}

func (s *serverImpl) Headers() map[string][]string {
	headers := make(map[string][]string, len(s.headers))
	for k, v := range s.headers {
		// Also copy value which is a slice
		headers[k] = append(make([]string, 0, len(v)), v...)
	}
	return headers
}

// addHeaders adds the extra headers to the specified response.
func (s *serverImpl) addHeaders(w http.ResponseWriter) {
	header := w.Header()
	for k, v := range s.headers {
		for _, v2 := range v {
			header.Add(k, v2)
		}
	}
}

func (s *serverImpl) AddStaticDir(path, dir string) error {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	if path == "" {
		return errors.New("path cannot be empty string")
	}

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	origPath := path
	path = s.appPath + path

	// pathEvent and pathRenderComp are window-relative so no need to check with those
	if path == s.appPath+pathStatic || path == s.appPath+pathSessCheck {
		return errors.New("Path cannot be '" + origPath + "' (reserved)!")
	}

	handler := http.StripPrefix(path, http.FileServer(http.Dir(dir)))
	// To include extra headers in the response of static handler:
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		s.addHeaders(w)
		handler.ServeHTTP(w, r)
	})

	return nil
}

func (s *serverImpl) Theme() string {
	return s.theme
}

func (s *serverImpl) SetTheme(theme string) {
	s.theme = theme
}

func (s *serverImpl) SetLogger(logger *log.Logger) {
	s.logger = logger
}

func (s *serverImpl) Logger() *log.Logger {
	return s.logger
}

func (s *serverImpl) AddRootHeadHTML(html string) {
	s.rootHeads = append(s.rootHeads, html)
}

func (s *serverImpl) RemoveRootHeadHTML(html string) {
	for i, v := range s.rootHeads {
		if v == html {
			old := s.rootHeads
			s.rootHeads = append(s.rootHeads[:i], s.rootHeads[i+1:]...)
			old[len(old)-1] = ""
			return
		}
	}
}

func (s *serverImpl) SetAppRootHandler(f AppRootHandlerFunc) {
	s.appRootHandlerFunc = f
}

func (s *serverImpl) SessIDCookieName() string {
	return s.sessIDCookieName
}

func (s *serverImpl) SetSessIDCookieName(name string) {
	s.sessIDCookieName = name
}

// serveStatic handles the static contents of server.
func (s *serverImpl) serveStatic(w http.ResponseWriter, r *http.Request) {
	s.addHeaders(w)

	// Parts example: "/appname/_gui_static/gui-0.8.0.js" => {"", "appname", "_gui_static", "gui-0.8.0.js"}
	parts := strings.Split(r.URL.Path, "/")

	if s.appName == "" {
		// No app name, gui server resides in root
		if len(parts) < 2 {
			// This should never happen. Path is always at least a slash ("/").
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string and pathStatic
		parts = parts[2:]
	} else {
		// We have app name
		if len(parts) < 3 {
			// Missing app name from path
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string, app name and pathStatic
		parts = parts[3:]
	}

	res := parts[0]
	if res == resNameStaticJs {
		w.Header().Set("Expires", time.Now().UTC().Add(72*time.Hour).Format(http.TimeFormat)) // Set 72 hours caching
		w.Header().Set("Content-Type", "application/x-javascript; charset=utf-8")
		w.Write(staticJs)
		return
	}
	if strings.HasSuffix(res, ".css") {
		cssCode := staticCSS[res]
		if cssCode != nil {
			w.Header().Set("Expires", time.Now().UTC().Add(72*time.Hour).Format(http.TimeFormat)) // Set 72 hours caching
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
			w.Write(cssCode)
			return
		}
	}

	http.NotFound(w, r)
}

// serveHTTP handles the incoming requests.
// Renders of the URL-selected window,
// and also handles event dispatching.
func (s *serverImpl) serveHTTP(w http.ResponseWriter, r *http.Request) {
	if s.logger != nil {
		s.logger.Println("Incoming:", r.URL.Path)
	}

	s.addHeaders(w)

	// Check session
	var sess Session
	c, err := r.Cookie(s.sessIDCookieName)
	if err == nil {
		s.sessMux.RLock()
		sess = s.sessions[c.Value]
		s.sessMux.RUnlock()
	}
	if sess == nil {
		sess = &s.sessionImpl
	}

	// Parts example: "/appname/winname/e?et=0&cid=1" => {"", "appname", "winname", "e"}
	parts := strings.Split(r.URL.Path, "/")

	if s.appName == "" {
		// No app name, gui server resides in root
		if len(parts) < 1 {
			// This should never happen. Path is always at least a slash ("/").
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string
		parts = parts[1:]
	} else {
		// We have app name
		if len(parts) < 2 {
			// Missing app name from path
			http.NotFound(w, r)
			return
		}
		// Omit the first empty string and the app name
		parts = parts[2:]
	}

	if len(parts) >= 1 && parts[0] == pathSessCheck {
		// Session check. Must not call sess.access()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		remaining := sess.Timeout() - time.Now().Sub(sess.Accessed())
		fmt.Fprintf(w, "%f", remaining.Seconds())
		return
	}

	if len(parts) < 1 || parts[0] == "" {
		// Missing window name, render window list
		s.appRootHandlerFunc(w, r, sess)
		return
	}

	winName := parts[0]

	win := sess.WinByName(winName)
	// If not found and we're on an authenticated session, try the public window list
	if win == nil && sess.Private() {
		win = s.WinByName(winName) // Server is a Session, the public session
		if win != nil {
			// We're serving a public window, switch to public session here entirely
			sess = &s.sessionImpl
		}
	}

	// If still not found and no private session, try the session creator names
	if win == nil && !sess.Private() {
		if _, found := s.sessCreatorNames[winName]; found {
			sess = s.newSession(nil)
			s.addSessCookie(sess, w)
			// Search again in the new session as SessionHandlers may have added windows.
			win = sess.WinByName(winName)
		}
	}

	if win == nil {
		// Invalid window name, render an error message with a link to the window list
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		NewWriter(w).Writess("<html><body>Window for name <b>'", winName, `'</b> not found. See the <a href="`, s.appPath, `">Window list</a>.</body></html>`)
		return
	}

	sess.access()

	var path string
	if len(parts) >= 2 {
		path = parts[1]
	}

	rwMutex := sess.rwMutex()
	switch path {
	case pathEvent:
		rwMutex.Lock()
		defer rwMutex.Unlock()

		s.handleEvent(sess, win, w, r)
	case pathRenderComp:
		rwMutex.RLock()
		defer rwMutex.RUnlock()

		// Render just a component
		s.renderComp(win, w, r)
	default:
		rwMutex.RLock()
		defer rwMutex.RUnlock()

		// Render the whole window
		win.RenderWin(NewWriter(w), s)
	}
}

// renderWinList builds a temporary Window, adds links to the windows of
// a session, and renders the Window.
func (s *serverImpl) renderWinList(wr http.ResponseWriter, r *http.Request, sess Session) {
	if s.logger != nil {
		s.logger.Println("\tRendering windows list.")
	}
	win := NewWindow("windowList", s.text+" - Window List")

	titleLabel := NewLabel(s.text + " - Window List")
	titleLabel.Style().SetFontWeight(FontWeightBold).SetFontSize("1.3em")
	win.Add(titleLabel)

	addLinks := func(title string, nameTexts [][2]string) {
		if len(nameTexts) == 0 {
			return
		}
		win.AddVSpace(10)
		win.Add(NewLabel(title))
		for _, nameText := range nameTexts {
			link := NewLink(nameText[1], path.Join(s.appPath, nameText[0]))
			link.Style().SetPaddingLeftPx(20)
			win.Add(link)
		}
	}

	// Render both private and public session windows
	sessions := make([]Session, 1, 2)
	sessions[0] = sess
	nameTexts := make([][2]string, 0, len(s.sessCreatorNames)+1)
	if sess.Private() {
		sessions = append(sessions, &s.sessionImpl)
	} else if len(s.sessCreatorNames) > 0 {
		// No private session yet, render session creators:
		nameTexts = nameTexts[:0]
		for name, text := range s.sessCreatorNames {
			nameTexts = append(nameTexts, [2]string{name, text})
		}
		addLinks("Session creators:", nameTexts)
	}

	for _, session := range sessions {
		text := "Public windows:"
		if session.Private() {
			text = "Authenticated windows:"
		}
		nameTexts = nameTexts[:0]
		for _, win := range session.SortedWins() {
			nameTexts = append(nameTexts, [2]string{win.Name(), win.Text()})
		}
		addLinks(text, nameTexts)
	}

	win.RenderWin(NewWriter(wr), s)
}

// renderComp renders just a component.
func (s *serverImpl) renderComp(win Window, w http.ResponseWriter, r *http.Request) {
	id, err := AtoID(r.FormValue(paramCompID))
	if err != nil {
		http.Error(w, "Invalid component id!", http.StatusBadRequest)
		return
	}

	if s.logger != nil {
		s.logger.Println("\tRendering comp:", id)
	}

	comp := win.ByID(id)
	if comp == nil {
		http.Error(w, fmt.Sprint("Component not found: ", id), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // We send it as text!
	comp.Render(NewWriter(w))
}

// handleEvent handles the event dispatching.
func (s *serverImpl) handleEvent(sess Session, win Window, wr http.ResponseWriter, r *http.Request) {
	focCompID, err := AtoID(r.FormValue(paramFocusedCompID))
	if err == nil {
		win.SetFocusedCompID(focCompID)
	}

	id, err := AtoID(r.FormValue(paramCompID))
	if err != nil {
		http.Error(wr, "Invalid component id!", http.StatusBadRequest)
		return
	}

	comp := win.ByID(id)
	if comp == nil {
		if s.logger != nil {
			s.logger.Println("\tComp not found:", id)
		}
		http.Error(wr, fmt.Sprint("Component not found: ", id), http.StatusBadRequest)
		return
	}

	etype := parseIntParam(r, paramEventType)
	if etype < 0 {
		http.Error(wr, "Invalid event type!", http.StatusBadRequest)
		return
	}
	if s.logger != nil {
		s.logger.Println("\tEvent from comp:", id, " event:", etype)
	}

	event := newEventImpl(EventType(etype), comp, s, sess, wr, r)
	shared := event.shared

	event.x = parseIntParam(r, paramMouseX)
	if event.x >= 0 {
		event.y = parseIntParam(r, paramMouseY)
		shared.wx = parseIntParam(r, paramMouseWX)
		shared.wy = parseIntParam(r, paramMouseWY)
		shared.mbtn = MouseBtn(parseIntParam(r, paramMouseBtn))
	} else {
		event.y, shared.wx, shared.wy, shared.mbtn = -1, -1, -1, -1
	}

	shared.modKeys = parseIntParam(r, paramModKeys)
	shared.keyCode = Key(parseIntParam(r, paramKeyCode))

	comp.preprocessEvent(event, r)

	// Dispatch event...
	comp.dispatchEvent(event)

	// Check if a new session was created during event dispatching
	if shared.session.New() {
		s.addSessCookie(shared.session, wr)
	}

	// ...and send back the result
	wr.Header().Set("Content-Type", "text/plain; charset=utf-8") // We send it as text
	w := NewWriter(wr)
	hasAction := false
	// If we reload, nothing else matters
	if shared.reload {
		hasAction = true
		w.Writevs(eraReloadWin, strComma, shared.reloadWin)
	} else {
		if len(shared.dirtyComps) > 0 {
			hasAction = true
			w.Writev(eraDirtyComps)
			for id := range shared.dirtyComps {
				w.Write(strComma)
				w.Writev(int(id))
			}
		}
		if shared.focusedComp != nil {
			if hasAction {
				w.Write(strSemicol)
			} else {
				hasAction = true
			}
			w.Writevs(eraFocusComp, strComma, int(shared.focusedComp.ID()))
			// Also register focusable comp at window
			win.SetFocusedCompID(shared.focusedComp.ID())
		}
	}
	if !hasAction {
		w.Writev(eraNoAction)
	}
}

// parseIntParam parses an int param.
// If error occurs, -1 will be returned.
func parseIntParam(r *http.Request, paramName string) int {
	if num, err := strconv.Atoi(r.FormValue(paramName)); err == nil {
		return num
	}
	return -1
}
