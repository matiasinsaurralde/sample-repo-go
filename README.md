# sample-repo-go

A small Go HTTP API built with [Gin](https://github.com/gin-gonic/gin).

Hello

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

### `GET /ls`

Lists the contents of a directory given by the required `path` query parameter.

**Request**

```bash
curl "http://localhost:8080/ls?path=/tmp"
```

**Response** `200 OK`

```json
{"output": "alpha.txt\nbeta.txt"}
```

Returns `400 Bad Request` if `path` is missing, and `500 Internal Server Error`
if the path cannot be listed. `path` is passed to `ls` as a single argument, so
shell metacharacters and wildcards are treated as part of the filename rather
than being expanded.
