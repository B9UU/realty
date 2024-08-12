package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
			realtyData: []*data.RealtyResponse{&mocks.MockRealties[0]},
			statusCode: http.StatusMethodNotAllowed,
		},
	}
	for _, test := range tests {

		app := newTestApp(test.realtyData)
		ts := newTestServer(t, app.routes())
		defer ts.Close()
		req := httptest.NewRequest(test.method, ts.URL+"/realty", nil)
		fmt.Println(req.RequestURI)
		req.RequestURI = ""
		resp, err := ts.Client().Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if status := resp.StatusCode; status != test.statusCode {
			t.Errorf("Hanlder got wrong http status want %d got %d", test.statusCode, status)
		}
		var gotBody map[string][]data.RealtyResponse
		if err := json.NewDecoder(resp.Body).Decode(&gotBody); err != nil {
			t.Fatal(err)
		}
		for _, i := range gotBody["realties"] {
			ddd := test.realtyData[0].ID
			if ddd != i.ID {
				t.Fatalf("not equal, want %v got %v", ddd, i.ID)
			}
		}
	}
}
func TestAddRealty(t *testing.T) {

}
