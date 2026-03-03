package server

// mongo and stuff
// redis

type Server struct {
}

func InitialiseServer() (*Server, error) {
	var s Server
	return &s, nil
}
