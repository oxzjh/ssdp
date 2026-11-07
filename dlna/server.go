package dlna

import (
	"net"
	"net/http"
)

type Server struct {
	l    net.Listener
	url  string
	file string
}

func (s *Server) GetURL() string {
	return s.url
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, s.file)
}

func (s *Server) Close() {
	s.l.Close()
}

func (s *Server) Serve() error {
	return http.Serve(s.l, s)
}

func NewServer(localIP, file string) (*Server, error) {
	l, err := net.Listen("tcp4", localIP+":0")
	if err != nil {
		return nil, err
	}
	return &Server{l, "http://" + l.Addr().String(), file}, nil
}
