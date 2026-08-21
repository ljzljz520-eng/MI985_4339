package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"example.com/nursery-cms/domain"
	"example.com/nursery-cms/service"
)

type HTTP struct{ Service *service.Service }

func NewHTTP(s *service.Service) *HTTP { return &HTTP{Service: s} }

func (h *HTTP) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/records", h.records)
	mux.HandleFunc("/import", h.importBatch)
	return mux
}

func (h *HTTP) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *HTTP) records(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := h.Service.Query(r.Context(), service.Query{BatchID: r.URL.Query().Get("batch"), Text: r.URL.Query().Get("q")})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func (h *HTTP) importBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload struct {
		BatchID string                  `json:"batch_id"`
		Items   []domain.ImportedRecord `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := h.Service.ImportBatch(r.Context(), payload.BatchID, payload.Items)
	if err != nil && !result.Cancelled {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	status := http.StatusOK
	if result.Cancelled {
		status = http.StatusRequestTimeout
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(result)
}

func parseStatus(value string) domain.RecordStatus {
	switch strings.ToLower(value) {
	case "draft":
		return domain.StatusDraft
	case "review":
		return domain.StatusReview
	case "approved":
		return domain.StatusApproved
	case "published":
		return domain.StatusPublished
	case "archived":
		return domain.StatusArchived
	default:
		return ""
	}
}

func (h *HTTP) Search(ctx context.Context, batch, status, text string) ([]domain.Record, error) {
	return h.Service.Query(ctx, service.Query{BatchID: batch, Status: parseStatus(status), Text: text})
}
