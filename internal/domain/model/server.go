package model

// Server represents a VPN edge server in the cluster.
type Server struct {
	ID         int64
	Code       string
	Domain     string
	IP         string
	Country    string
	UsersCount int64
	UsersMax   int64
}

// IsFull returns true if current users count reached maximum capacity.
func (s *Server) IsFull() bool {
	return s.UsersCount >= s.UsersMax
}
