package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"velozient-backend/api"
	"velozient-backend/db"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetPasswordCards(t *testing.T) {
	testcases := []struct {
		name           string
		dbload         []db.PasswordCard
		expectedStatus int
	}{
		{
			name: "TEST GET ALL SUCCESS",
			dbload: []db.PasswordCard{
				{URL: "url-1", UserName: "username-1", Name: "name-1", Password: "psswd-1"},
				{URL: "url-2", UserName: "username-2", Name: "name-2", Password: "psswd-2"},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "TEST GET ALL EMPTY REPLY",
			dbload:         []db.PasswordCard{},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			database := db.NewMemoryDB()
			database.PopulateDB(tc.dbload)

			expectedResponse := database.GetAllCards()
			r := httptest.NewRequest(http.MethodGet, "https://test.com/password-card", nil)
			w := httptest.NewRecorder()

			server := api.NewServer("8080", database)

			server.CreateOrGetPasswordCards(w, r)

			response := w.Result()

			body, _ := io.ReadAll(response.Body)

			var actual []db.PasswordCard
			_ = json.Unmarshal(body, &actual)

			assert.Equal(t, tc.expectedStatus, response.StatusCode)
			assert.ElementsMatch(t, expectedResponse, actual)
		})
	}
}

func TestCreatePasswordCards(t *testing.T) {
	testcases := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedUUID   string
	}{
		{
			name:           "TEST CREATE SUCCESS",
			payload:        `{"uuid":"","url": "url", "userName":"user", "name":"name", "password":"pwd"}`,
			expectedStatus: http.StatusOK,
			expectedUUID:   "uuid",
		},
		{
			name:           "TEST CREATE MISSING FIELD",
			payload:        `{"password":"pwd"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedUUID:   "",
		},
		{
			name:           "TEST CREATE INVALID PAYLOAD",
			payload:        `{"uuid":"","url": "url", "userName":"user", "name":"name", "password":"pwd","INVALIDFIELD":"INVALID"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedUUID:   "",
		},
	}

	for _, tc := range testcases {
		dbmock := &db.DBMock{}
		dbmock.On("CreateCard", mock.Anything).Return(tc.expectedUUID)
		server := api.NewServer("8080", dbmock)

		r := httptest.NewRequest(http.MethodPost,
			"https://test.com/password-card",
			bytes.NewReader([]byte(tc.payload)))

		w := httptest.NewRecorder()
		server.CreateOrGetPasswordCards(w, r)

		response := w.Result()
		body, _ := io.ReadAll(response.Body)

		var cardresponse api.CardCreateResponse
		_ = json.Unmarshal(body, &cardresponse)

		assert.Equal(t, tc.expectedStatus, response.StatusCode)
		assert.Equal(t, tc.expectedUUID, cardresponse.UUID)
	}
}

func TestEditCard(t *testing.T) {
	testcases := []struct {
		name           string
		payload        string
		expectedStatus int
		cardNotFound   error
	}{
		{
			name:           "TEST EDIT SUCCESS",
			payload:        `{"uuid":"noop","url": "url", "userName":"user", "name":"name", "password":"pwd"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "TEST EDIT FAIL CARD NOT FOUND",
			payload:        `{"uuid":"noop","url": "url", "userName":"user", "name":"name", "password":"pwd"}`,
			expectedStatus: http.StatusBadRequest,
			cardNotFound:   errors.New("card not found"),
		},
		{
			name:           "TEST CREATE MISSING FIELD",
			payload:        `{"password":"pwd"}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "TEST CREATE INVALID PAYLOAD",
			payload:        `{"uuid":"","url": "url", "userName":"user", "name":"name", "password":"pwd","INVALIDFIELD":"INVALID"}`,
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tc := range testcases {
		dbmock := &db.DBMock{}
		dbmock.On("EditCard", mock.Anything).Return(tc.cardNotFound)
		server := api.NewServer("8080", dbmock)

		r := httptest.NewRequest(http.MethodPut,
			"https://test.com/password-cards/uuid",
			bytes.NewReader([]byte(tc.payload)))

		w := httptest.NewRecorder()
		server.EditOrDeletePasswordCards(w, r)

		response := w.Result()

		assert.Equal(t, tc.expectedStatus, response.StatusCode)
	}
}

func TestDeleteCard(t *testing.T) {
	// This is always ok
	database := db.NewMemoryDB()
	server := api.NewServer("8080", database)
	expectedStatus := http.StatusOK
	r := httptest.NewRequest(http.MethodDelete, "https://test.com/password-card/", nil)
	w := httptest.NewRecorder()

	server.EditOrDeletePasswordCards(w, r)
	response := w.Result()

	assert.Equal(t, expectedStatus, response.StatusCode)
}
