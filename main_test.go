// main_test.go

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHostnamesHandler(t *testing.T) {
	// define test cases for threshold from 1 to 10
	testCases := []struct {
		threshold string
		expected  []string
	}{
		{"1", []string{"mta-prod-1"}},
		{"2", []string{"mta-prod-1", "mta-prod-2"}},
		{"3", []string{"mta-prod-1", "mta-prod-2"}},
		{"4", []string{"mta-prod-1", "mta-prod-2"}},
		{"5", []string{"mta-prod-1", "mta-prod-2"}},
		{"6", []string{"mta-prod-1", "mta-prod-2"}},
		{"7", []string{"mta-prod-1", "mta-prod-2"}},
		{"8", []string{"mta-prod-1", "mta-prod-2"}},
		{"9", []string{"mta-prod-1", "mta-prod-2"}},
		{"10", []string{"mta-prod-1", "mta-prod-2"}},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("GET", "/hostnames?x="+string(tc.threshold), nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(hostnamesHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var result []string
		err = json.NewDecoder(rr.Body).Decode(&result)
		if err != nil {
			t.Fatal(err)
		}

		// Check if expected hostnames are present in the result
		for _, expectedHostname := range tc.expected {
			if !stringInSlice(expectedHostname, result) {
				t.Errorf("handler returned unexpected result: got %v, expected %v", result, tc.expected)
			}
		}
	}
}

func stringInSlice(s string, slice []string) bool {
	for _, elem := range slice {
		if elem == s {
			return true
		}
	}
	return false
}
