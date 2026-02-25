# google-integrity-api
HTTP server to handle Google Play Integrity API requests

## Usage
### Start the server
`go run ./cmd/server/`

### Generate a new nonce
```bash
$ curl localhost:8080/nonce/generate
{"nonce":"NuaA7Aw3zKpUERwk0lnrKioCfNGnfvdmSUFOMtbzP0s="}
```

### Verify
```bash
$ curl -X POST http://localhost:8080/nonce/verify \
  -H "Content-Type: application/json" \
  -d '{"nonce":"NuaA7Aw3zKpUERwk0lnrKioCfNGnfvdmSUFOMtbzP0s="}'
{"valid":true}
```