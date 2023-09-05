package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"velozient-backend/db"
)

const base10 = 10

// Server holds the necessary
// attributes of our http server.
type Server struct {
	mux  *http.ServeMux
	port string
	db   *db.MemoryDB
}

// NewServer returns a routless server with a db connection.
func NewServer(port string, db *db.MemoryDB) *Server {
	mux := http.NewServeMux()
	return &Server{
		db:   db,
		mux:  mux,
		port: port,
	}
}

// CreateOrGetPasswordCards handles requests for the getting all
// cards or creating a new card through the /password-cards route.
func (s *Server) CreateOrGetPasswordCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet:
			fmt.Println("GET /password-cards")
			w.WriteHeader(http.StatusOK)

			cards := s.db.GetAllCards()
			fmt.Println(cards)
			data, err := json.Marshal(cards)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if _, err := w.Write(data); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		case r.Method == http.MethodPost:
			fmt.Println("POST /password-cards")
			w.WriteHeader(http.StatusOK)
		default:
			fmt.Println("Not allowedfound")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

// EditOrDeletePasswordCards handles requests for editing or deleting
// a card identified by :id in the /password-cards/:id route.
func (s *Server) EditOrDeletePasswordCards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idParam := strings.TrimPrefix(r.URL.Path, "/password-cards/")
		id, err := strconv.ParseUint(idParam, base10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch {
		case r.Method == http.MethodPut:
			fmt.Printf("PUT /password-cards/%d\n", id)
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodDelete:
			fmt.Printf("DELETE /password-cards/%d\n", id)
			w.WriteHeader(http.StatusOK)
		default:
			fmt.Println("Not allowed")
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
}

// RegisterRoutes centralizes the list
// of routes and respective  handlers
// for the current server.
func (s *Server) RegisterRoutes() {
	// Could've used a regex based pattern matching implementation
	// for this route, however, to keep things simple, the current
	// http.ServeMux approach was preferred.
	s.mux.Handle("/password-cards", s.CreateOrGetPasswordCards())
	s.mux.Handle("/password-cards/", s.EditOrDeletePasswordCards())
}

// Run starts our http server
func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.port)
	return http.ListenAndServe(address, s.mux)
}
