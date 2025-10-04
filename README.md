# mozzy — Postman‑level JSON HTTP client for the terminal

`mozzy` is a modern, developer‑friendly alternative to `curl`, designed for JSON APIs and token‑heavy auth workflows. Think **Postman on the CLI**: clean verbs, JSON pretty‑print, inline querying, request history, JWT decode/verify/sign, environments, and request **chaining** via captures and `{{vars}}`.

## Highlights
- Simple verbs: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`
- JSON pretty‑print + `--jq` inline querying
- Auth & headers: `--auth $TOKEN`, `--header 'X-Env: staging'`
- Request history & replay
- Built‑in JWT decode / verify / sign
- Environments via `.mozzy.json`
- **Chaining:** `--capture name=.path`, use `{{name}}` later
- **Workflows:** `mozzy run flow.yaml` (YAML steps with captures and assertions)

## Install (local build)
```bash
cd starter-boilerplate
go build -o mozzy ./...
./mozzy --help
```

## Quick start
```bash
./mozzy --base https://api.example.com GET /users
./mozzy POST /login --json '{"user":"alice","pass":"secret"}'
./mozzy GET /users --jq '.[0].id'
./mozzy jwt decode <token>
```

## Chaining (captures + vars)
```bash
./mozzy POST /auth --json @creds.json --capture token=.access_token
./mozzy GET /profile --auth "{{token}}"
```

## Workflow YAML
See [examples/user-flow.yaml](examples/user-flow.yaml) then run:
```bash
./mozzy run examples/user-flow.yaml --env staging
```

## Roadmap
- v1: verbs, JSON, JWT, history, capture/vars, run YAML
- v1.1: collections, assertions, retries, cookie jar, `--as-curl`, HAR export
- v2: diff/snapshot, benchmark, mock server, OpenAPI import, OAuth device flow

## License
MIT © 2025 Archith Sharma
