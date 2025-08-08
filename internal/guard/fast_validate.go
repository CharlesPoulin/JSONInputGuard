package guard

import (
    "bytes"
    "encoding/json"
    "errors"
    "unicode"

    "github.com/example/jsoninputguard/internal/types"
    "github.com/example/jsoninputguard/internal/validate"
)

// predictGuardShape parses top-level fields without decoding the heavy features array.
type predictGuardShape struct {
	UserID    string          `json:"user_id" validate:"required,min=1,max=64"`
	SessionID string          `json:"session_id" validate:"required,min=1,max=64"`
	Timestamp int64           `json:"timestamp" validate:"required"`
	Features  json.RawMessage `json:"features" validate:"required"`
	Metadata  map[string]string `json:"metadata" validate:"max=128,dive,keys,max=64,endkeys,max=4096"`
}

// fastValidatePredict validates structure and sizes in ~O(n) over raw bytes.
func fastValidatePredict(buf []byte) error {
	    var g predictGuardShape
    if err := json.Unmarshal(buf, &g); err != nil {
        return err
    }
	// Validate light fields using validator singleton
	if err := validate.V().Struct(&g); err != nil {
		return err
	}
	// Validate features as array length within bounds without decoding numbers
	if len(g.Features) == 0 {
		return errors.New("features: empty")
	}
	if !isJSONArray(g.Features) {
		return errors.New("features: not array")
	}
	count, ok := fastCountArrayItems(g.Features)
	if !ok {
		return errors.New("features: invalid array syntax")
	}
	if count < 1 || count > 16384 {
		return errors.New("features: length out of bounds")
	}
	return nil
}

func isJSONArray(b []byte) bool {
	i := 0
	for i < len(b) && unicode.IsSpace(rune(b[i])) { i++ }
	if i >= len(b) || b[i] != '[' { return false }
	j := len(b) - 1
	for j >= 0 && unicode.IsSpace(rune(b[j])) { j-- }
	return j >= 0 && b[j] == ']'
}

// fastCountArrayItems counts items in a JSON array without full parse.
// It handles nested whitespace and numbers/strings, assuming elements themselves do not contain nested arrays/objects.
// This is sufficient for numeric feature arrays.
func fastCountArrayItems(b []byte) (int, bool) {
	// Trim outer brackets
	l := bytes.IndexByte(b, '[')
	if l < 0 { return 0, false }
	r := bytes.LastIndexByte(b, ']')
	if r < 0 || r <= l { return 0, false }
	inner := b[l+1:r]
	// Empty array
	skip := func(p int) int {
		for p < len(inner) && (inner[p] == ' ' || inner[p] == '\n' || inner[p] == '\t' || inner[p] == '\r') { p++ }
		return p
	}
	p := skip(0)
	if p >= len(inner) { return 0, true }
	count := 1
	inString := false
	for p < len(inner) {
		c := inner[p]
		if c == '"' {
			inString = !inString
			p++
			continue
		}
		if !inString && c == ',' {
			count++
		}
		p++
	}
	return count, true
}

// validateAndMaybeDecode runs the fast validator and optionally decodes into full struct.
func validateAndMaybeDecode[T any](buf []byte, dst *T) error {
	switch any(dst).(type) {
	case *types.PredictRequest:
		if err := fastValidatePredict(buf); err != nil {
			return err
		}
	}
	// Decode into destination
	    return json.Unmarshal(buf, dst)
}

// GuardAndDecodePredict validates PredictRequest using sonic AST for minimal work,
// then decodes once into the destination struct.
// GuardPredictRaw performs fast structural and size validation without decoding the heavy array.
func GuardPredictRaw(buf []byte) error {
    // Single-pass, zero-allocation scanner for top-level fields
    const (
        stateKey = iota
        stateColon
        stateValue
    )

    i := 0
    // Skip leading spaces
    for i < len(buf) && (buf[i] == ' ' || buf[i] == '\n' || buf[i] == '\r' || buf[i] == '\t') { i++ }
    if i >= len(buf) || buf[i] != '{' {
        return errors.New("invalid json: not object")
    }
    i++

    var haveUser, haveSess, haveTS, haveFeat bool
    var userLen, sessLen int
    var tsOK bool
    var featCount int

    for i < len(buf) {
        // Skip whitespace and commas
        for i < len(buf) && (buf[i] == ' ' || buf[i] == '\n' || buf[i] == '\r' || buf[i] == '\t' || buf[i] == ',') { i++ }
        if i >= len(buf) {
            break
        }
        if buf[i] == '}' {
            i++
            break
        }

        // Parse key string
        if buf[i] != '"' {
            return errors.New("invalid json: key not string")
        }
        i++
        keyStart := i
        escaped := false
        for i < len(buf) {
            c := buf[i]
            if c == '"' && !escaped {
                break
            }
            if c == '\\' && !escaped {
                escaped = true
            } else {
                escaped = false
            }
            i++
        }
        if i >= len(buf) {
            return errors.New("invalid json: unterminated key")
        }
        key := buf[keyStart:i]
        i++ // skip closing quote

        // Skip to ':'
        for i < len(buf) && (buf[i] == ' ' || buf[i] == '\n' || buf[i] == '\r' || buf[i] == '\t') { i++ }
        if i >= len(buf) || buf[i] != ':' {
            return errors.New("invalid json: missing colon")
        }
        i++
        for i < len(buf) && (buf[i] == ' ' || buf[i] == '\n' || buf[i] == '\r' || buf[i] == '\t') { i++ }
        if i >= len(buf) {
            return errors.New("invalid json: missing value")
        }

        // Match keys we care about and validate value
        switch {
        case bytes.Equal(key, []byte("user_id")):
            if buf[i] != '"' { return errors.New("user_id: not string") }
            i++
            start := i
            escaped = false
            for i < len(buf) {
                c := buf[i]
                if c == '"' && !escaped { break }
                if c == '\\' && !escaped { escaped = true } else { escaped = false }
                i++
            }
            if i >= len(buf) { return errors.New("user_id: unterminated") }
            userLen = i - start
            if userLen < 1 || userLen > 64 { return errors.New("user_id: length out of bounds") }
            i++ // closing quote
            haveUser = true

        case bytes.Equal(key, []byte("session_id")):
            if buf[i] != '"' { return errors.New("session_id: not string") }
            i++
            start := i
            escaped = false
            for i < len(buf) {
                c := buf[i]
                if c == '"' && !escaped { break }
                if c == '\\' && !escaped { escaped = true } else { escaped = false }
                i++
            }
            if i >= len(buf) { return errors.New("session_id: unterminated") }
            sessLen = i - start
            if sessLen < 1 || sessLen > 64 { return errors.New("session_id: length out of bounds") }
            i++ // closing quote
            haveSess = true

        case bytes.Equal(key, []byte("timestamp")):
            // parse optional minus and digits
            sign := int64(1)
            if buf[i] == '-' { sign = -1; i++ }
            if i >= len(buf) || buf[i] < '0' || buf[i] > '9' { return errors.New("timestamp: invalid") }
            var n int64
            for i < len(buf) && buf[i] >= '0' && buf[i] <= '9' {
                n = n*10 + int64(buf[i]-'0')
                i++
            }
            n *= sign
            tsOK = n > 0
            if !tsOK { return errors.New("timestamp: must be > 0") }
            haveTS = true

        case bytes.Equal(key, []byte("features")):
            if buf[i] != '[' { return errors.New("features: not array") }
            c, end, ok := countArrayItemsAndEnd(buf, i)
            if !ok { return errors.New("features: invalid array") }
            if c < 1 || c > 16384 { return errors.New("features: length out of bounds") }
            featCount = c
            haveFeat = true
            i = end

        default:
            // Skip unknown value (string/number/obj/array) to next comma/brace
            // Basic skipper for performance: handles strings and bracketed structures
            switch buf[i] {
            case '"':
                i++
                escaped = false
                for i < len(buf) {
                    c := buf[i]
                    if c == '"' && !escaped { i++; break }
                    if c == '\\' && !escaped { escaped = true } else { escaped = false }
                    i++
                }
            case '{':
                depth := 1; i++
                for i < len(buf) && depth > 0 {
                    if buf[i] == '"' {
                        i++
                        escaped = false
                        for i < len(buf) {
                            c := buf[i]
                            if c == '"' && !escaped { i++; break }
                            if c == '\\' && !escaped { escaped = true } else { escaped = false }
                            i++
                        }
                        continue
                    }
                    if buf[i] == '{' { depth++ }
                    if buf[i] == '}' { depth-- }
                    i++
                }
            case '[':
                _, end, ok := countArrayItemsAndEnd(buf, i)
                if !ok { return errors.New("invalid array") }
                i = end
            default:
                // number, true, false, null
                for i < len(buf) && buf[i] != ',' && buf[i] != '}' && buf[i] != '\n' && buf[i] != '\r' && buf[i] != '\t' && buf[i] != ' ' { i++ }
            }
        }
        // After value, continue loop to next key or end
    }

    if !haveUser || !haveSess || !haveTS || !haveFeat || featCount == 0 {
        return errors.New("missing required fields")
    }
    return nil
}

// countArrayItemsAndEnd counts items in JSON array starting at '[' and returns (count, endIndexAfterBracket, ok)
func countArrayItemsAndEnd(buf []byte, start int) (int, int, bool) {
    i := start
    if buf[i] != '[' { return 0, i, false }
    i++
    // Skip whitespace
    for i < len(buf) && (buf[i] == ' ' || buf[i] == '\n' || buf[i] == '\r' || buf[i] == '\t') { i++ }
    if i < len(buf) && buf[i] == ']' {
        return 0, i+1, true
    }
    count := 1
    inString := false
    escaped := false
    depth := 1
    for i < len(buf) {
        c := buf[i]
        if inString {
            if c == '"' && !escaped { inString = false }
            if c == '\\' && !escaped { escaped = true } else { escaped = false }
            i++
            continue
        }
        switch c {
        case '"':
            inString = true
        case '[':
            // do not support nested arrays for features; treat as invalid
            return 0, i, false
        case ']':
            depth--
            if depth == 0 { return count, i+1, true }
        case ',':
            count++
        }
        i++
    }
    return 0, i, false
}

// findJSONStringValue finds a string value for a given field name at top level.
func findJSONStringValue(buf []byte, field string) (string, bool) {
    // naive, robust enough for top-level search
    key := []byte("\"" + field + "\"")
    i := bytes.Index(buf, key)
    if i < 0 {
        return "", false
    }
    // move to colon after key
    j := i + len(key)
    for j < len(buf) && (buf[j] == ' ' || buf[j] == '\t' || buf[j] == '\n' || buf[j] == '\r') {
        j++
    }
    if j >= len(buf) || buf[j] != ':' {
        return "", false
    }
    j++
    for j < len(buf) && (buf[j] == ' ' || buf[j] == '\t' || buf[j] == '\n' || buf[j] == '\r') {
        j++
    }
    if j >= len(buf) || buf[j] != '"' {
        return "", false
    }
    j++
    start := j
    escaped := false
    for j < len(buf) {
        c := buf[j]
        if c == '"' && !escaped {
            return string(buf[start:j]), true
        }
        if c == '\\' && !escaped {
            escaped = true
        } else {
            escaped = false
        }
        j++
    }
    return "", false
}

// findJSONNumberValue finds a positive integer value for a given field name.
func findJSONNumberValue(buf []byte, field string) (int64, bool) {
    key := []byte("\"" + field + "\"")
    i := bytes.Index(buf, key)
    if i < 0 { return 0, false }
    j := i + len(key)
    for j < len(buf) && (buf[j] == ' ' || buf[j] == '\t' || buf[j] == '\n' || buf[j] == '\r') { j++ }
    if j >= len(buf) || buf[j] != ':' { return 0, false }
    j++
    for j < len(buf) && (buf[j] == ' ' || buf[j] == '\t' || buf[j] == '\n' || buf[j] == '\r') { j++ }
    var n int64
    if j >= len(buf) { return 0, false }
    sign := int64(1)
    if buf[j] == '-' { sign = -1; j++ }
    if j >= len(buf) || buf[j] < '0' || buf[j] > '9' { return 0, false }
    for j < len(buf) && buf[j] >= '0' && buf[j] <= '9' {
        n = n*10 + int64(buf[j]-'0')
        j++
    }
    return n*sign, true
}

// findJSONArrayRaw returns the raw bytes of the array value for a given field.
func findJSONArrayRaw(buf []byte, field string) ([]byte, bool) {
    key := []byte("\"" + field + "\"")
    i := bytes.Index(buf, key)
    if i < 0 { return nil, false }
    j := i + len(key)
    for j < len(buf) && (buf[j] == ' ' || buf[j] == '\t' || buf[j] == '\n' || buf[j] == '\r') { j++ }
    if j >= len(buf) || buf[j] != ':' { return nil, false }
    j++
    for j < len(buf) && (buf[j] == ' ' || buf[j] == '\t' || buf[j] == '\n' || buf[j] == '\r') { j++ }
    if j >= len(buf) || buf[j] != '[' { return nil, false }
    start := j
    depth := 0
    inString := false
    escaped := false
    for j < len(buf) {
        c := buf[j]
        if inString {
            if c == '"' && !escaped { inString = false }
            if c == '\\' && !escaped { escaped = true } else { escaped = false }
            j++
            continue
        }
        switch c {
        case '"':
            inString = true
        case '[':
            depth++
        case ']':
            depth--
            if depth == 0 {
                return buf[start : j+1], true
            }
        }
        j++
    }
    return nil, false
}

// GuardAndDecodePredict validates then decodes.
func GuardAndDecodePredict(buf []byte, dst *types.PredictRequest) error {
    if err := GuardPredictRaw(buf); err != nil {
        return err
    }
    return json.Unmarshal(buf, dst)
}
