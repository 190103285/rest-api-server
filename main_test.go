package main

import (
	"log"
	"os"
	"testing"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		"postgres",
		"password",
		"coffeeshop",
		"disable")
	//,
	//os.Getenv("disable")

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM coffee")
	a.DB.Exec("ALTER SEQUENCE id RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS coffee
(
    id BIGSERIAL,
    name TEXT NOT NULL,
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/coffees", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentCoffee(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/coffee/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Coffee not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Coffee not found'. Got '%s'", m["error"])
	}
}

func TestCreateCoffee(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"name":"test coffee"}`)
	req, _ := http.NewRequest("POST", "/coffee", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test coffee" {
		t.Errorf("Expected coffee name to be 'test coffee'. Got '%v'", m["name"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected coffee ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetCoffee(t *testing.T) {
	clearTable()
	addCoffee(1)

	req, _ := http.NewRequest("GET", "/coffee/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCoffee(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO coffee(name) VALUES($1)")
	}
}

func TestUpdateCoffee(t *testing.T) {

	clearTable()
	addCoffee(1)

	req, _ := http.NewRequest("GET", "/coffee/1", nil)
	response := executeRequest(req)
	var originalCoffee map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCoffee)

	var jsonStr = []byte(`{"name":"test coffee - updated name"}`)
	req, _ = http.NewRequest("PUT", "/coffee/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalCoffee["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalCoffee["id"], m["id"])
	}

	if m["name"] == originalCoffee["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalCoffee["name"], m["name"], m["name"])
	}
}

func TestDeleteCoffee(t *testing.T) {
	clearTable()
	addCoffee(1)

	req, _ := http.NewRequest("GET", "/coffee/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/coffee/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/coffee/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
