package kingologs

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

// Server is an interface describing kingologs server objects.
type Server struct {
	logger     Logger
	TargetChan chan string
	config     ConfigValues
}

// Get the server started.
func (s Server) start() {
	connStr := fmt.Sprintf("%s:%d", s.config.Connection.TCP.Host, s.config.Connection.TCP.Port)
	s.logger.Info.Printf("Starting server at: %s", connStr)
	listener, err := net.Listen("tcp", connStr)
	if err != nil {
		s.logger.Error.Printf("Starting Server error: %s", err)
	}

	for {
		conn, err := listener.Accept()
		// TODO: handle errors with one client gracefully.
		if err != nil {

			s.logger.Info.Println("Accepting Connection...")
		}

		go s.handleRequest(conn)
	}
}

func (s Server) handleRequest(conn net.Conn) {
	s.logger.Trace.Println("Handling new connection request")
	for {
		buf := make([]byte, 512)

		_, err := conn.Read(buf)
		if err != nil {
			s.logger.Error.Printf("Reading Connection error: %s", err)
			return
		}
		defer conn.Close()

		// Get rid of newline and bull chars on buffer.
		tmp := strings.TrimSpace(string(bytes.Trim(buf, "\x00")))

		s.logger.Trace.Printf("Got from client: %#v", tmp)
		s.TargetChan <- tmp
		return
	}
}

// NewServer gets us a new server.
//func NewServer(l Logger, c ConfigValues) *Server {
func NewServer(l Logger, c ConfigValues) *Server {
	l.Info.Println("Getting ready to create a new server")

	// Start the new server and prep it.
	s := new(Server)
	s.logger = l
	s.config = c

	l.Trace.Println("Done creating a server")

	return s
}

// StartServer gets the server started.
func (s *Server) StartServer() {
	s.logger.Info.Println("Beginning to start server")
	go s.start()
	s.logger.Trace.Println("Done starting server")
}

// SetTargetChan sets the target channel.
func (s *Server) SetTargetChan(t chan string) {
	s.TargetChan = t
}
