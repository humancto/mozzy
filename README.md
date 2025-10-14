# mozzy — Postman for your Terminal 🚀

`mozzy` is a modern, developer-friendly HTTP client built for JSON APIs and token-heavy workflows. Think **Postman on the CLI**: colorized output, inline querying, request collections, JWT tools, workflows, and powerful request **chaining** via captures and variables.

## ✨ Features

### Core HTTP
- **Clean verbs**: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`
- **Beautiful output**: Auto-colorized JSON with TTY detection
- **Inline queries**: `--jq` for JSONPath filtering (nested paths, arrays)
- **Response times**: Automatic timing for every request
- **Smart errors**: Helpful HTTP status explanations with actionable tips
- **Verbose mode**: `--verbose` for headers & timing breakdown
- **Retry logic**: `--retry N` with exponential backoff
- **Cookie jar**: `--cookie-jar` for session persistence

### Authentication & Headers
- **Bearer auth**: `--auth <token>`
- **Custom headers**: `--header 'X-API-Key: secret'` (repeatable)
- **Environments**: Named configs in `.mozzy.json`

### Request Management
- **Collections**: Save, list, and execute requests
  ```bash
  mozzy save api-login POST /auth --json @creds.json
  mozzy list
  mozzy exec api-login
  ```
- **History**: Browse recent requests with colors
  ```bash
  mozzy history --limit 20
  ```

### Chaining & Workflows
- **Variable captures**: Extract values from responses
  ```bash
  mozzy GET /users --capture userId=[0].id
  mozzy GET /users/{{userId}}/posts
  ```
- **YAML workflows**: Multi-step API flows with captures
  ```bash
  mozzy run workflow.yaml
  ```

### JWT Tools
- **Decode**: `mozzy jwt decode <token>`
- **Verify**: `mozzy jwt verify <token> --secret <key>`
- **Sign**: `mozzy jwt sign payload.json --secret <key>`
- **JWKS support**: `--jwk <url>` for RSA/ECDSA verification

## 🚀 Quick Start

### Installation

#### Homebrew (macOS & Linux) - RECOMMENDED
```bash
brew tap humancto/mozzy
brew install mozzy
```

#### From Source
```bash
git clone https://github.com/humancto/mozzy.git
cd mozzy
go build -o mozzy .

# Optional: Install globally
sudo cp mozzy /usr/local/bin/
```

#### Pre-built Binaries
Download from [GitHub Releases](https://github.com/humancto/mozzy/releases)

### Try it out
```bash
# Simple GET with colors
mozzy GET https://jsonplaceholder.typicode.com/users/1

# Filter with --jq
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city

# POST with JSON
mozzy POST https://api.example.com/login --json '{"user":"alice"}'

# Save to collection
mozzy save my-api GET https://api.example.com/users
mozzy exec my-api

# View history
mozzy history --limit 10
```

## 📚 Examples

### API Chaining
```bash
# Login and capture token
mozzy POST /auth --json @creds.json --capture token=.access_token

# Use token in subsequent requests
mozzy GET /profile --auth "{{token}}"
mozzy GET /data --auth "{{token}}" --jq .results[0]
```

### Workflows
Create `flow.yaml`:
```yaml
name: User API Flow
steps:
  - name: Get user list
    method: GET
    url: https://api.example.com/users
    capture:
      userId: .[0].id

  - name: Get user details
    method: GET
    url: https://api.example.com/users/{{userId}}
```

Run it:
```bash
mozzy run flow.yaml
```

### Environments
Create `.mozzy.json`:
```json
{
  "environments": {
    "dev": {
      "base_url": "http://localhost:3000",
      "headers": {"X-Env": "development"}
    },
    "prod": {
      "base_url": "https://api.example.com",
      "auth_token": "your-token"
    }
  }
}
```

Use it:
```bash
mozzy --env prod GET /users
mozzy env  # List available environments
```

## Roadmap
- v1: verbs, JSON, JWT, history, capture/vars, run YAML
- v1.1: collections, assertions, retries, cookie jar, `--as-curl`, HAR export
- v2: diff/snapshot, benchmark, mock server, OpenAPI import, OAuth device flow

## License
MIT © 2025 Archith Sharma
