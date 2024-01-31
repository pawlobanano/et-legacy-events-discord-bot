package api

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

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
	fs := http.FileServer(http.Dir("./static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", s.serveTemplate)
	mux.HandleFunc("/cup", s.getTeamLineupsByDefaultEdition)
	mux.HandleFunc("/cup/edition", s.getTeamLineupsByEditionID)
	mux.HandleFunc("/cup/team", s.getTeamLineupByDefaultEditionByTeamIDLetter)

	return http.ListenAndServe(s.listenAddr, mux)
}

// serveTemplate function builds paths to the layout file and the template file corresponding with the request.
// Finally, the template.ExecuteTemplate() function renders a named template in the set, in this case the layout template.
func (s *Server) serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

	// Return a 404 if the template doesn't exist.
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory.
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		s.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
