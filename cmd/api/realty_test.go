package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/b9uu/realty/internal/data"
	"github.com/b9uu/realty/internal/mocks"
)

func TestGetRealties(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		realtyData []*data.RealtyResponse
		statusCode int
	}{
		{
			name:       "with bad method",
			method:     http.MethodPut,
			realtyData: nil,
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "Valid method",
			method:     http.MethodGet,
			realtyData: []*data.RealtyResponse{&mocks.MockRealties[0], &mocks.MockRealties[1]},
			statusCode: http.StatusOK,
		},
	}
	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			app := newTestApp()
			ts := newTestServer(t, app.routes())
			defer ts.Close()
			// initiate new http request
			req := httptest.NewRequest(test.method, ts.URL+"/realty", nil)
			// httptest.NewRequest setup RequestURI since it's meant for a server to send the request not a client
			req.RequestURI = ""
			// call the server client to do the request
			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check for status code
			if status := resp.StatusCode; status != test.statusCode {
				t.Errorf("Hanlder got wrong http status want %d got %d", test.statusCode, status)
			}
			// FIX: not doing anything useful
			if len(test.realtyData) > 0 {
				var gotBody map[string][]data.RealtyResponse
				if err := json.NewDecoder(resp.Body).Decode(&gotBody); err != nil {
					t.Fatal(err)
				}
				for i, v := range gotBody["realties"] {
					ddd := test.realtyData[i].ID
					if ddd != v.ID {
						t.Fatalf("not equal, want %v got %v", ddd, v.ID)
					}
				}
			}
		})
	}
}
func TestAutoComplete(t *testing.T) {
	var ValidInput = "Van"
	var InvalidInput = "V"
	var MethodNotAllowed = "Method Not Allowed"

	tests := []struct {
		name       string
		method     string
		statusCode int
		want       string
		input      string
	}{
		{
			name:       "with bad method",
			method:     http.MethodPut,
			statusCode: http.StatusMethodNotAllowed,
			want:       MethodNotAllowed,
			input:      ValidInput,
		},
		{
			name:       "Valid method",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			want:       "result",
			input:      ValidInput,
		},
		{
			name:       "With bad input",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			want:       "error",
			input:      InvalidInput,
		},
	}

	for _, test := range tests {

		app := newTestApp()
		ts := newTestServer(t, app.routes())
		defer ts.Close()
		// initiate new http request
		req := httptest.NewRequest(test.method, ts.URL+"/auto-complete", nil)
		// add query
		q := req.URL.Query()
		q.Add("city", test.input)

		req.URL.RawQuery = q.Encode()
		// httptest.NewRequest setup RequestURI since it's meant for a server to send the request not a client
		req.RequestURI = ""
		// call the server client to do the request
		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// check for status code
		if status := resp.StatusCode; status != test.statusCode {
			t.Errorf("Hanlder got wrong http status want %d got %d", test.statusCode, status)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), test.want) {
			t.Errorf("Body do not contain %s, got %s", test.want, string(body))
		}

	}
}
