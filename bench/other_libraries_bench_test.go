package bench

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fastjson"
	"github.com/tidwall/gjson"
	"github.com/example/jsoninputguard/internal/guard"
	"github.com/example/jsoninputguard/internal/types"
	"github.com/example/jsoninputguard/internal/validate"
)

// Benchmark JSONInputGuard (reusing existing function)
// func BenchmarkGuardOnly64KB(b *testing.B) - already exists in guard_only_bench_test.go

// Benchmark standard library (reusing existing function)  
// func BenchmarkStdlibUnmarshalOnly64KB(b *testing.B) - already exists in guard_only_bench_test.go

// Benchmark jsoniter
func BenchmarkJsoniterUnmarshal64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	var sink any
	for i := 0; i < b.N; i++ {
		var v any
		if err := jsoniter.Unmarshal(payload, &v); err != nil {
			b.Fatalf("err: %v", err)
		}
		sink = v
	}
	_ = sink
}

// Benchmark fastjson
func BenchmarkFastjsonParse64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	var sink *fastjson.Value
	for i := 0; i < b.N; i++ {
		v, err := fastjson.ParseBytes(payload)
		if err != nil {
			b.Fatalf("err: %v", err)
		}
		sink = v
	}
	_ = sink
}

// Benchmark gjson
func BenchmarkGjsonParse64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	var sink gjson.Result
	for i := 0; i < b.N; i++ {
		result := gjson.ParseBytes(payload)
		sink = result
	}
	_ = sink
}

// Benchmark with validation for each library
func BenchmarkGuardWithValidation64KB(b *testing.B) {
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

func BenchmarkJsoniterWithValidation64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	for i := 0; i < b.N; i++ {
		var req types.PredictRequest
		if err := jsoniter.Unmarshal(payload, &req); err != nil {
			b.Fatalf("err: %v", err)
		}
		if err := validate.V().Struct(&req); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}

func BenchmarkStdlibWithValidation64KB(b *testing.B) {
	payload := makePayload(64 * 1024)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	for i := 0; i < b.N; i++ {
		var req types.PredictRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			b.Fatalf("err: %v", err)
		}
		if err := validate.V().Struct(&req); err != nil {
			b.Fatalf("err: %v", err)
		}
	}
}

// Memory usage comparison for different payload sizes
func BenchmarkMemoryComparison(b *testing.B) {
	sizes := []int{1024, 4096, 16384, 65536} // 1KB, 4KB, 16KB, 64KB
	
	for _, size := range sizes {
		payload := makePayload(size)
		b.Run(fmt.Sprintf("Guard_%dKB", size/1024), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(payload)))
			for i := 0; i < b.N; i++ {
				if err := guard.GuardPredictRaw(payload); err != nil {
					b.Fatalf("err: %v", err)
				}
			}
		})
		
		b.Run(fmt.Sprintf("Stdlib_%dKB", size/1024), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(payload)))
			for i := 0; i < b.N; i++ {
				var v any
				if err := json.Unmarshal(payload, &v); err != nil {
					b.Fatalf("err: %v", err)
				}
			}
		})
		
		b.Run(fmt.Sprintf("Jsoniter_%dKB", size/1024), func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(payload)))
			for i := 0; i < b.N; i++ {
				var v any
				if err := jsoniter.Unmarshal(payload, &v); err != nil {
					b.Fatalf("err: %v", err)
				}
			}
		})
	}
}
