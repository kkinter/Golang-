package server

import "newspjt/handlers"

func (s *Server) setupRoutes() {
	handlers.Health(s.mux)
}
