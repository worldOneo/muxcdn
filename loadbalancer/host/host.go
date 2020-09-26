package host

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Server and everything it needs
type Server struct {
	Name         string
	URL          string
	Status       *ServerStatus
	ReverseProxy *httputil.ReverseProxy
}

// ServerStatus the status of the server
type ServerStatus struct {
	Online bool
}

// New creates a new server
func New(Name string, URL string) (*Server, error) {
	parsedURL, err := url.Parse(URL)

	if err != nil {
		return nil, err
	}
	return &Server{
		Name: Name,
		URL:  URL,
		Status: &ServerStatus{
			Online: false,
		},
		ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}, nil
}

// Handle lets the server handle the reques/response
func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	s.ReverseProxy.ServeHTTP(w, r)
}

// IsRunning pings the server and returns if the server can be reached
func (s *Server) IsRunning() (bool, error) {
	resp, err := http.Head(s.URL)
	if err != nil || resp == nil {
		return false, err
	}
	if resp.StatusCode >= 500 {
		return false, nil
	}
	return true, nil
}
