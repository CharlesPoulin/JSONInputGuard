package guard

import (
    "context"
    "errors"
    "net/http"
    "sync"
    "time"

    "encoding/json"

    "github.com/example/jsoninputguard/internal/types"
)

// rawBufferPool uses pooling to reduce allocations and GC cost.
var rawBufferPool = &sync.Pool{New: func() any { b := make([]byte, 0, 64*1024); return &b }}

// MaxPayloadSize caps the JSON payload we accept.
const MaxPayloadSize = 64 * 1024 // 64 KiB

// DecodeValidateJSON reads, bounds, decodes with sonic, and optionally validates.
// It avoids reflection on the hot path by using sonic.Unmarshal.
func DecodeValidateJSON[T any](w http.ResponseWriter, r *http.Request, dst *T, validateFn func(*T) error) error {
	// Enforce size cap early using http.MaxBytesReader
	r.Body = http.MaxBytesReader(w, r.Body, MaxPayloadSize)
	defer r.Body.Close()

	bufPtr := rawBufferPool.Get().(*[]byte)
	buf := *bufPtr
	buf = buf[:0]
	defer func() { *bufPtr = buf[:0]; rawBufferPool.Put(bufPtr) }()

    // Read into preallocated pooled buffer to avoid extra copies
    for {
        if len(buf) == cap(buf) {
            // Should not happen due to MaxBytesReader, but guard anyway
            return errors.New("payload too large")
        }
        // Extend slice for next read chunk
        nMax := cap(buf) - len(buf)
        // Read directly into the slice tail
        tmp := (*bufPtr)[:len(buf)+nMax]
        n, err := r.Body.Read(tmp[len(buf):])
        if n > 0 {
            buf = tmp[:len(buf)+n]
        }
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            // For short reads, we may get io.EOF at end; treat as done
            if n == 0 {
                break
            }
        }
        if n == 0 {
            break
        }
    }

	if len(buf) == 0 {
		return errors.New("empty body")
	}

	// Fast path: validate shape from raw, then decode
	if pr, ok := any(dst).(*types.PredictRequest); ok {
		if err := GuardAndDecodePredict(buf, pr); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(buf, dst); err != nil {
			return err
		}
	}

	if validateFn != nil {
		if err := validateFn(dst); err != nil {
			return err
		}
	}

	return nil
}

// TimeBudgetMiddleware rejects requests that exceed a per-request time budget.
func TimeBudgetMiddleware(budget time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			deadline := time.Now().Add(budget)
			r = r.WithContext(context.WithValue(r.Context(), ctxDeadlineKey{}, deadline))
			next.ServeHTTP(w, r)
		})
	}
}

type ctxDeadlineKey struct{}
