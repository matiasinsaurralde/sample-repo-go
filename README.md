# sample-repo-go

A small Go HTTP API built with [Gin](https://github.com/gin-gonic/gin).

## Project layout

```
.
├── cmd/server/     # Application entrypoint
├── pkg/api/        # HTTP handlers and routing
└── pkg/config/     # Configuration
```

## Requirements

- Go 1.25+

## Run

```bash
go run ./cmd/server
```

The server listens on `:8080` by default. Override with the `ADDR` environment variable:

```bash
ADDR=:3000 go run ./cmd/server
```

## API

### `POST /hello`

Accepts a JSON body with an optional `name` field and returns a greeting.

**Request**

```bash
curl -X POST http://localhost:8080/hello \
  -H "Content-Type: application/json" \
  -d '{"name": "world"}'
```

**Response** `200 OK`

```json
{"message": "hello world"}
```

If `name` is omitted, it defaults to `"world"`.
