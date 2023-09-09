package db

import (
	"fmt"
	"sync"

	"github.com/rs/xid"
)

// PasswordCard defines the mandatory
// fields for a password card.
type PasswordCard struct {
	UUID     string `json:"uuid"`
	URL      string `json:"url"`
	UserName string `json:"userName"`
	Name     string `json:"name"`
	Password string `json:"password"` // This wouldn't be plain text IRL (:
}

// MemoryDB encapsulates the needed
// methods for the applications data layer.
type MemoryDB interface {
	CreateCard(card PasswordCard) string
	GetAllCards() []PasswordCard
	EditCard(card PasswordCard) error
	DeleteCard(uuid string)
}

// memDB holds access to the
// applications in-memory data base.
type memDB struct {
	mux   sync.Mutex // needed because of the concurrent nature of the http calls
	cards map[string]PasswordCard
}

// NewMemoryDB creates an
// empty in-memory data base.
func NewMemoryDB() *memDB {
	cards := make(map[string]PasswordCard)
	return &memDB{
		cards: cards,
	}
}

// PopulateDB exists with the sole purpose
// of populating the application's database
// for a good presentation, since this is
// a technical assignment.
func (m *memDB) PopulateDB(load []PasswordCard) {
	for _, card := range load {
		guid := xid.New()
		uuid := guid.String()
		card.UUID = uuid
		m.cards[uuid] = card
	}
}

// CreateCard creates a new card with given info
// and stores it in the in-memory database and
// returns the uuid of the newly created password card.
func (m *memDB) CreateCard(card PasswordCard) string {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Generate unique identifier
	guid := xid.New()
	uuid := guid.String()

	card.UUID = uuid
	m.cards[uuid] = card

	return uuid
}

// GetAllCards retrieves all stored password cards.
func (m *memDB) GetAllCards() []PasswordCard {
	m.mux.Lock()
	defer m.mux.Unlock()

	cards := make([]PasswordCard, 0)
	for _, card := range m.cards {
		cards = append(cards, card)
	}
	return cards
}

// EditCard edits a given card and returns an
// error in case the password card doesn't exist.
func (m *memDB) EditCard(card PasswordCard) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	_, ok := m.cards[card.UUID]
	if !ok {
		return fmt.Errorf("card not found: %s", card.UUID)
	}

	m.cards[card.UUID] = card
	return nil
}

// DeleteCard deletes a the card identified by uuid.
// In case uuid is not found, DeleteCards is a no-op.
func (m *memDB) DeleteCard(uuid string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.cards, uuid)
}
