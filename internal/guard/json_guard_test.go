package guard

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPayload struct {
	Name  string `json:"name" validate:"required"`
	Value int    `json:"value" validate:"gt=0"`
}

func TestDecodeValidateJSON_Valid(t *testing.T) {
	payload := testPayload{Name: "test", Value: 1}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	var result testPayload
	err := DecodeValidateJSON(rr, req, &result, func(p *testPayload) error { return nil })

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, payload, result)
}

func TestDecodeValidateJSON_InvalidJSON(t *testing.T) {
	body := []byte(`{"name": "test", "value": 1,}`)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	var result testPayload
	err := DecodeValidateJSON(rr, req, &result, func(p *testPayload) error { return nil })

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDecodeValidateJSON_ValidationError(t *testing.T) {
	payload := testPayload{Name: "", Value: 1}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	var result testPayload
	err := DecodeValidateJSON(rr, req, &result, func(p *testPayload) error { return errors.New("validation error") })

	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
