package bench

import (
	"testing"

	"encoding/json"

	"github.com/example/jsoninputguard/internal/guard"
)

func BenchmarkGuardOnly64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	for i := 0; i < b.N; i++ {
		if err := guard.GuardPredictRaw(payload); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}

func BenchmarkStdlibUnmarshalOnly64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	var sink any
	for i := 0; i < b.N; i++ {
		var v any
		if err := json.Unmarshal(payload, &v); err != nil {
			b.Fatalf("err: %v", err)
		}
		sink = v
	}
	_ = sink
}
