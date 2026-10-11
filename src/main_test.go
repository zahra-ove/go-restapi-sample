package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zahra-ove/go-restapi-sample/db"
	"github.com/zahra-ove/go-restapi-sample/models"
)

func TestEventPersistenceWithLegacySchema(t *testing.T) {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	conn.SetMaxOpenConns(1)
	previous := db.DB
	db.DB = conn
	defer func() { db.DB = previous }()
	_, err = conn.Exec(`CREATE TABLE events (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, description TEXT NOT NULL, location TEXT NOT NULL, datatime DATETIME NOT NULL, user_id INTEGER)`)
	if err != nil {
		t.Fatal(err)
	}
	db.CreateTables()
	db.CreateTables() // Repeated startup must also work.
	router := gin.New()
	router.POST("/events", CreateEvent)
	router.GET("/events", getEvents)
	for i := 0; i < 2; i++ {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"Name":"Persistence test","Description":"Test","Location":"Test","DateTime":"2026-10-11T12:00:00Z"}`))
		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(response, request)
		if response.Code != http.StatusCreated {
			t.Fatalf("POST: %d %s", response.Code, response.Body.String())
		}
		var body struct{ Event models.Event }
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body.Event.ID != int64(i+1) {
			t.Fatalf("unexpected ID: %d", body.Event.ID)
		}
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/events", nil))
	var events []models.Event
	if err := json.Unmarshal(response.Body.Bytes(), &events); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK || len(events) != 2 {
		t.Fatalf("GET: %d %s", response.Code, response.Body.String())
	}
	if _, err := conn.Exec(`DROP TABLE events`); err != nil {
		t.Fatal(err)
	}
	response = httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(`{"Name":"Test","Description":"Test","Location":"Test","DateTime":"2026-10-11T12:00:00Z"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("failed save returned %d", response.Code)
	}
}
