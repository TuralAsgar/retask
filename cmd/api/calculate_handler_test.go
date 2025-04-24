package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/TuralAsgar/dynamic-programming/internal/data"
	"github.com/julienschmidt/httprouter"
	_ "modernc.org/sqlite"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", "test.sqlite")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE pack_size (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			size INTEGER NOT NULL UNIQUE
		)
	`)
	if err != nil {
		t.Fatalf("failed to create pack_size table: %v", err)
	}

	return db
}

func seedTestDB(t *testing.T, db *sql.DB, packSizes []int) {
	t.Helper()

	for _, size := range packSizes {
		_, err := db.Exec(`INSERT INTO pack_size (size) VALUES (?)`, size)
		if err != nil {
			t.Fatalf("failed to seed pack_size table: %v", err)
		}
	}
}

func teardownTestDB(t *testing.T, db *sql.DB) {
	t.Helper()

	err := db.Close()
	if err != nil {
		t.Fatalf("failed to close test database: %v", err)
	}

	err = os.Remove("test.sqlite")
	if err != nil {
		t.Fatalf("failed to remove test database file: %v", err)
	}
}

func TestCalculateHandler(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	seedTestDB(t, db, []int{250, 500, 1000, 2000, 5000})

	app := &application{
		models: data.Models{
			Calculator: &data.CalculatorModel{DB: db},
		},
	}

	testCases := []struct {
		name             string
		orderAmount      int
		expectedPackages map[int]int
		expectedStatus   int
	}{
		{
			name:        "Order 1 item",
			orderAmount: 1,
			expectedPackages: map[int]int{
				250: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Order exactly 250 items",
			orderAmount: 250,
			expectedPackages: map[int]int{
				250: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Order 251 items",
			orderAmount: 251,
			expectedPackages: map[int]int{
				500: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Order 501 items",
			orderAmount: 501,
			expectedPackages: map[int]int{
				500: 1,
				250: 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Order 12001 items",
			orderAmount: 12001,
			expectedPackages: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/pack/calculate/%d", tc.orderAmount), nil)

			// Add the httprouter.Params to the request context
			// Because we use https://github.com/julienschmidt/httprouter
			// And it does not get "size" param from the request url directly
			params := httprouter.Params{httprouter.Param{Key: "size", Value: strconv.Itoa(tc.orderAmount)}}
			ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			app.calculatePackSizeHandler(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			var actualResponse map[string]map[int]int
			err := json.NewDecoder(rr.Body).Decode(&actualResponse)
			if err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if len(actualResponse["packages"]) != len(tc.expectedPackages) {
				t.Fatalf("expected %d package types; got %d. Expected: %v, Got: %v",
					len(tc.expectedPackages),
					len(actualResponse["packages"]),
					tc.expectedPackages,
					actualResponse["packages"])
			}

			for size, count := range tc.expectedPackages {
				if actualResponse["packages"][size] != count {
					t.Errorf("expected %d packs of size %d; got %d", count, size, actualResponse["packages"][size])
				}
			}
		})
	}
}
