# Google Play Integrity API
HTTP server to handle Google Play Integrity API requests

## Usage
### Start the server
`go run ./cmd/server/`

### Generate new one-time token
```bash
$ curl localhost:8080/token/generate
{"token":"NuaA7Aw3zKpUERwk0lnrKioCfNGnfvdmSUFOMtbzP0s="}
```

### Verify
```bash
$ curl -X POST http://localhost:8080/token/verify \
  -H "Content-Type: application/json" \
  -d '{"token":"NuaA7Aw3zKpUERwk0lnrKioCfNGnfvdmSUFOMtbzP0s="}'
{"valid":true}
```