package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/b9uu/realty/internal/data"
)

type realtyResp struct {
	Realties []data.Realty `json:"realties"`
}

func TestRealty(t *testing.T) {
	app := newTestApp()
	req, err := http.NewRequest(http.MethodGet, "/realty", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.getRealties)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Hanlder got wrong http status want %d got %d", http.StatusOK, status)
	}
	var got realtyResp
	err = json.NewDecoder(rr.Body).Decode(&got)
	if err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
	for i, r := range got.Realties {
		if r.ID != i+1 {
			t.Errorf("want %d got %d: %+v", i+1, r.ID, r)
		}
		t.Logf("want %d got %d: %+v", i+1, r.ID, r)

	}

}
