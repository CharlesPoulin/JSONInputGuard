package bench

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/example/jsoninputguard/internal/predict"
)

func BenchmarkHTTPPredict64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	h := predict.Router()
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/predict", bytes.NewReader(payload))
		h.ServeHTTP(w, r)
		if w.Code != 200 {
			b.Fatalf("status: %d", w.Code)
		}
	}
}
