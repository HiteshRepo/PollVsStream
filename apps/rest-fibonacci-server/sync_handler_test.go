package restfibonacciserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func newTestRouter() *mux.Router {
	app := NewApp()
	router := mux.NewRouter()
	router.HandleFunc("/fibonacci/sync/{number:[0-9]+}", app.fibonacciSyncHandler)
	return router
}

func TestSyncHandlerValidInput(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/fibonacci/sync/10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp SyncResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.FibonacciNumbers) != 10 {
		t.Errorf("expected 10 fibonacci numbers, got %d", len(resp.FibonacciNumbers))
	}
}

func TestSyncHandlerZeroInput(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/fibonacci/sync/0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp SyncResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.FibonacciNumbers) != 0 {
		t.Errorf("expected 0 fibonacci numbers for n=0, got %d", len(resp.FibonacciNumbers))
	}
}

func TestSyncHandlerContentType(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/fibonacci/sync/5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
}
