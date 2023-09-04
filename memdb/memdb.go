package memdb

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
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Password string `json:"password"` // This wouldn't be plain text IRL (:
}

// MemoryDB holds access to the
// applications in-memory data base.
type MemoryDB struct {
	mux   *sync.Mutex // needed because of the concurrent nature of the http calls
	cards map[string]PasswordCard
}

// NewMemoryDB creates an
// empty in-memory data base.
func NewMemoryDB() *MemoryDB {
	cards := make(map[string]PasswordCard)
	return &MemoryDB{
		cards: cards,
	}
}

// CreateCard creates a new card with given info
// and stores it in the in-memory database and
// returns the uuid of the newly created password card.
func (m *MemoryDB) CreateCard(card PasswordCard) string {
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
func (m *MemoryDB) GetAllCards() []PasswordCard {
	m.mux.Lock()
	defer m.mux.Unlock()

	return m.cards
}

// EditCard edits a given card and returns an
// error in case the password card doesn't exist.
func (m *MemoryDB) EditCard(card PasswordCard) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	card, ok := m.cards[card.UUID]
	if !ok {
		return fmt.Errorf("card not found: %s", card.UUID)
	}

	m.cards[card.UUID] = card

	return nil
}

// DeleteCard deletes a the card identified by uuid.
// In case uuid is not found, DeleteCards is a no-op.
func (m *MemoryDB) DeleteCard(uuid string) {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.cards, uuid)
	return nil
}
