package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"velozient-backend/db"
)

// Server holds the necessary
// attributes of our http server.
type Server struct {
	handler  http.Handler
	port     string
	database db.MemoryDB
}

// CardCreateResponse holds the
// structured reply for a card
// creationg POST request.
type CardCreateResponse struct {
	UUID string `json:"uuid"`
}

// NewServer returns a routless server with a db connection.
func NewServer(port string, database db.MemoryDB) *Server {
	return &Server{
		database: database,
		port:     port,
	}
}

// CreateOrGetPasswordCards handles requests for the getting all
// cards or creating a new card through the /password-cards route.
func (s *Server) CreateOrGetPasswordCards(w http.ResponseWriter, r *http.Request) {
	switch {

	case r.Method == http.MethodGet:
		w.WriteHeader(http.StatusOK)

		cards := s.database.GetAllCards()
		database, err := json.Marshal(cards)
		if err != nil {
			fmt.Printf("error reading from database: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(database); err != nil {
			fmt.Printf("error writing response body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	case r.Method == http.MethodPost:
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r.Body)
		if err != nil {
			fmt.Printf("error reading request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := validate(buf.Bytes()); err != nil {
			fmt.Printf("invalid request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var card db.PasswordCard
		if err := json.Unmarshal(buf.Bytes(), &card); err != nil {
			fmt.Printf("error unmarshaling request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		uuid := s.database.CreateCard(card)

		database, err := json.Marshal(&CardCreateResponse{
			UUID: uuid,
		})
		if err != nil {
			fmt.Printf("error marshaling response body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(database); err != nil {
			fmt.Printf("error writing response body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

	default:
		fmt.Println("Not allowedfound")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// EditOrDeletePasswordCards handles requests for editing or deleting
// a card identified by :id in the /password-cards/:id route.
func (s *Server) EditOrDeletePasswordCards(w http.ResponseWriter, r *http.Request) {
	idParam := strings.TrimPrefix(r.URL.Path, "/password-cards/")

	switch {
	case r.Method == http.MethodPut:
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r.Body)
		if err != nil {
			fmt.Printf("error reading request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := validate(buf.Bytes()); err != nil {
			fmt.Printf("invalid request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var card db.PasswordCard
		if err := json.Unmarshal(buf.Bytes(), &card); err != nil {
			fmt.Printf("error unmarshaling request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.database.EditCard(card); err != nil {
			fmt.Printf("error editing card: %s\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case r.Method == http.MethodDelete:
		s.database.DeleteCard(idParam)
		w.WriteHeader(http.StatusOK)

	default:
		fmt.Println("Not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return
}

// validate validates the input data
// coming from requests to the server.
// this would be more properly done through
// the use of reflection or a third party
// module like https://github.com/go-playground/validator.
// However, this simple implementation
// suffices as example for our simple server.
func validate(b []byte) error {
	// disallow unknown fields
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.DisallowUnknownFields()

	var passcard db.PasswordCard
	if err := decoder.Decode(&passcard); err != nil {
		return err
	}

	// check if there are enough fields
	payloadfields := make(map[string]any)
	if err := json.Unmarshal(b, &payloadfields); err != nil {
		return err
	}

	if len(payloadfields) != 5 {
		return errors.New("missing mandatory field")
	}

	return nil
}

// RegisterRoutes centralizes the list
// of routes and respective  handlers
// for the current server.
func (s *Server) Setup() {
	mux := http.NewServeMux()
	mux.HandleFunc("/password-cards", s.CreateOrGetPasswordCards)
	mux.HandleFunc("/password-cards/", s.EditOrDeletePasswordCards)

	s.handler = WithLoggerMiddleware(WithCORSMiddleware(mux))
}

// Run starts our http server
func (s *Server) Run() error {
	address := fmt.Sprintf(":%s", s.port)
	fmt.Println("Starting server...")
	return http.ListenAndServe(address, s.handler)
}
