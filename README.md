## Ultra-fast JSON input guard (chi + sonic + validator)

Targets:
- Validate ~64 kB payloads in ≤ 100 µs p95 for guard-only path (AST-based checks)
- Adds ≤ 1 ms to /predict handler end-to-end
- Cold start under 120 ms on AWS Lambda (Chi adapter)

Commands:
- Build: `make build`
- Run: `make run` (listens on :8080)
- Test: `make test`
- Benchmarks: `make bench`, `make bench-guard`, and `make benchmark-all`
- Lambda artifact: `make lambda` -> `lambda.zip`

Notes:
- Guard uses `sonic/ast` search to validate top-level fields and count array items without decoding (O(n) over raw bytes) and decodes once.
- `http.MaxBytesReader` caps payloads at 64 KiB.
- Minimal middleware to keep latency budget tight.

AWS Lambda:
- Uses `aws-lambda-go-api-proxy/chi` for API Gateway compatibility.
- Build with `GOOS=linux GOARCH=amd64 CGO_ENABLED=0` to ensure small, fast binaries.
