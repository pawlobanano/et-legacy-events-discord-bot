package api

import (
	"log/slog"
	"net/http"

	"github.com/pawlobanano/et-legacy-events-discord-bot/types"
)

// Server struct sets a type used by http package to run HTTP server.
type Server struct {
	config     *types.Environemnt
	logger     *slog.Logger
	listenAddr string
}

// NewServer function creates and returns a new Server struct.
func NewServer(config *types.Environemnt, logger *slog.Logger) *Server {
	return &Server{
		config:     config,
		listenAddr: config.SERVER_ADDRESS,
		logger:     logger,
	}
}

// Start function starts HTTP server on providedd address (set in Server struct).
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/cup", s.getTeamLineupsByDefaultEdition)
	mux.HandleFunc("/cup/edition", s.getTeamLineupsByEditionID)
	// mux.HandleFunc("/cup/team", s.getTeamLineupOfSetEditionByTeamID)

	return http.ListenAndServe(s.listenAddr, mux)
}
