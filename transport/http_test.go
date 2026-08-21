package transport

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/nursery-cms/service"
	"example.com/nursery-cms/store"
)

func TestHTTPHealthAndImport(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/http.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	s := service.New(db)
	h := NewHTTP(s).Handler()
	health := httptest.NewRecorder()
	h.ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/health", nil))
	if health.Code != http.StatusOK || !strings.Contains(health.Body.String(), "ok") {
		t.Fatalf("health=%d %s", health.Code, health.Body.String())
	}
	body := `{"batch_id":"batch-http","items":[{"external_id":"1","batch_id":"batch-http","title":"活动","content":"内容","owner":"teacher"}]}`
	result := httptest.NewRecorder()
	h.ServeHTTP(result, httptest.NewRequest(http.MethodPost, "/import", strings.NewReader(body)))
	if result.Code != http.StatusOK || !strings.Contains(result.Body.String(), "batch-http") {
		t.Fatalf("import=%d %s", result.Code, result.Body.String())
	}
}
