package bench

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/example/jsoninputguard/internal/guard"
	"github.com/example/jsoninputguard/internal/types"
	"github.com/example/jsoninputguard/internal/validate"
)

func makePayload(targetBytes int) []byte {
	features := make([]float32, 0, 4096)
	for len(features) < 8192 { // tune up to ~64 KiB JSON
		features = append(features, 0.1234, 1.2345, 2.3456, 3.4567)
		b, _ := json.Marshal(map[string]any{
			"user_id":    "u",
			"session_id": "s",
			"timestamp":  1,
			"features":   features,
		})
		if len(b) >= targetBytes {
			return b
		}
	}
	b, _ := json.Marshal(map[string]any{
		"user_id":    "u",
		"session_id": "s",
		"timestamp":  1,
		"features":   features,
	})
	return b
}

func BenchmarkDecodeValidate1KB(b *testing.B) {
	payload := makePayload(1 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		var req types.PredictRequest
		r := httptest.NewRequest("POST", "/predict", bytes.NewReader(payload))
		if err := guard.DecodeValidateJSON(w, r, &req, func(p *types.PredictRequest) error { return validate.V().Struct(p) }); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}

func BenchmarkDecodeValidate64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		var req types.PredictRequest
		r := httptest.NewRequest("POST", "/predict", bytes.NewReader(payload))
		if err := guard.DecodeValidateJSON(w, r, &req, func(p *types.PredictRequest) error { return validate.V().Struct(p) }); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}

func BenchmarkDecodeValidate1MB(b *testing.B) {
	payload := makePayload(1 * 1024 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	w := httptest.NewRecorder()
	for i := 0; i < b.N; i++ {
		var req types.PredictRequest
		r := httptest.NewRequest("POST", "/predict", bytes.NewReader(payload))
		if err := guard.DecodeValidateJSON(w, r, &req, func(p *types.PredictRequest) error { return validate.V().Struct(p) }); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}
