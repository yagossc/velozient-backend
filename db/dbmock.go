package db

import "github.com/stretchr/testify/mock"

// DBMock implements the
// MemoryDB interface for
// mocking in tests.
type DBMock struct {
	mock.Mock
}

func (d *DBMock) CreateCard(card PasswordCard) string {
	args := d.Called(card)
	return args.String(0)
}
func (d *DBMock) GetAllCards() []PasswordCard {
	args := d.Called()
	return args.Get(0).([]PasswordCard)

}
func (d *DBMock) EditCard(card PasswordCard) error {
	args := d.Called(card)
	return args.Error(0)
}

func (d *DBMock) DeleteCard(uuid string) {
	// noop
}
