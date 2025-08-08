package predict

import (
	"net/http"
	"time"

	"encoding/json"
	"github.com/go-chi/chi/v5"

	"github.com/example/jsoninputguard/internal/guard"
	"github.com/example/jsoninputguard/internal/types"
	"github.com/example/jsoninputguard/internal/validate"
)

// Router returns a chi router with the /predict route.
func Router() *chi.Mux {
	r := chi.NewRouter()
	// Minimal middleware to keep latency budget tight. Add a soft time budget.
	r.Use(guard.TimeBudgetMiddleware(950 * time.Millisecond))

	r.Post("/predict", PredictHandler)
	return r
}

func PredictHandler(w http.ResponseWriter, r *http.Request) {
	var req types.PredictRequest
	if err := guard.DecodeValidateJSON(w, r, &req, func(p *types.PredictRequest) error {
		return validate.V().Struct(p)
	}); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	// Simulate a very cheap scoring function to isolate guard overhead.
	score := fastScore(req.Features)
	resp := types.PredictResponse{Score: score}
	writeJSON(w, http.StatusOK, resp)
}

func fastScore(features []float32) float32 {
	var s float32
	// Sum first 16 items at most to keep deterministic and cheap
	limit := 16
	if len(features) < limit {
		limit = len(features)
	}
	for i := 0; i < limit; i++ {
		s += features[i]
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Minimal overhead stdlib marshal
	b, _ := json.Marshal(v)
	_, _ = w.Write(b)
}
