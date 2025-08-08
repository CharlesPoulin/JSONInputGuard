package types

// PredictRequest is the incoming JSON payload for /predict.
// Fields chosen to exercise ~64 KB payloads at typical array sizes.
type PredictRequest struct {
	UserID     string    `json:"user_id" validate:"required,min=1,max=64"`
	SessionID  string    `json:"session_id" validate:"required,min=1,max=64"`
	Timestamp  int64     `json:"timestamp" validate:"required"`
	Features   []float32 `json:"features" validate:"required,min=1,max=16384,dive"`
	Metadata   map[string]string `json:"metadata" validate:"max=128,dive,keys,max=64,endkeys,max=4096"`
}

// PredictResponse is a compact response.
type PredictResponse struct {
	Score float32 `json:"score"`
}
