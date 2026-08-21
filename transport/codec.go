package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"example.com/nursery-cms/domain"
)

type ImportEnvelope struct {
	BatchID string                  `json:"batch_id"`
	Items   []domain.ImportedRecord `json:"items"`
}

func DecodeImport(reader io.Reader) (ImportEnvelope, error) {
	var envelope ImportEnvelope
	if reader == nil {
		return envelope, fmt.Errorf("import body is nil")
	}
	if err := json.NewDecoder(reader).Decode(&envelope); err != nil {
		return envelope, err
	}
	if envelope.BatchID == "" {
		return envelope, fmt.Errorf("batch_id is required")
	}
	return envelope, nil
}

func EncodeImport(value domain.ImportResult) ([]byte, error) { return json.Marshal(value) }

func EncodeRecords(records []domain.Record) ([]byte, error) { return json.Marshal(records) }

func ParseQuery(value string) map[string]string {
	result := make(map[string]string)
	for _, part := range strings.Split(value, "&") {
		pieces := strings.SplitN(part, "=", 2)
		if len(pieces) == 2 && pieces[0] != "" {
			result[pieces[0]] = pieces[1]
		}
	}
	return result
}

func ErrorPayload(err error) []byte {
	message := "unknown error"
	if err != nil {
		message = err.Error()
	}
	payload, marshalErr := json.Marshal(map[string]string{"error": message})
	if marshalErr != nil {
		return []byte(`{"error":"marshal"}`)
	}
	return payload
}

func StatusCode(err error) int {
	if err == nil {
		return 200
	}
	if strings.Contains(err.Error(), "required") {
		return 400
	}
	if strings.Contains(err.Error(), "conflict") {
		return 409
	}
	return 500
}
