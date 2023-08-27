package server

func New(options Options) *Server {
	return &Server{options: options}
}
