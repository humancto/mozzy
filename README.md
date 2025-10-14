# mozzy — Postman for your Terminal 🚀

[![GitHub release](https://img.shields.io/github/v/release/humancto/mozzy)](https://github.com/humancto/mozzy/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://go.dev/)

A modern, developer-friendly HTTP client built for the terminal. Think **Postman meets curl** with beautiful colors, inline JSON queries, request collections, JWT tools, and powerful workflow automation.

```bash
# One command to rule them all
mozzy GET https://api.github.com/users/torvalds --jq .name --color
```

---

## ✨ Why mozzy?

**Stop fighting with curl syntax.** mozzy gives you:

- 🎨 **Beautiful colors** - Auto-colorized JSON that's easy on the eyes
- 🔍 **Inline queries** - Filter JSON with `--jq` without piping to jq
- 📚 **Collections** - Save and reuse requests like Postman
- 🔗 **API chaining** - Capture values and use them in next requests
- ⚙️ **Workflows** - Multi-step API flows in YAML
- 🔐 **JWT superpowers** - Decode, verify, sign JWTs instantly
- 🚀 **Dev-friendly** - Built by developers, for developers

---

## 🚀 Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/humancto/mozzy/main/install.sh | bash
```

### Homebrew

```bash
brew tap humancto/mozzy
brew install mozzy
```

### From Binary

Download for your platform from [GitHub Releases](https://github.com/humancto/mozzy/releases/latest):

- macOS (Intel): `mozzy_darwin_amd64.tar.gz`
- macOS (Apple Silicon): `mozzy_darwin_arm64.tar.gz`
- Linux (x64): `mozzy_linux_amd64.tar.gz`
- Linux (ARM): `mozzy_linux_arm64.tar.gz`
- Windows: `mozzy_windows_amd64.zip`

```bash
# Extract and install
tar -xzf mozzy_*.tar.gz
sudo mv mozzy /usr/local/bin/
```

### From Source

```bash
git clone https://github.com/humancto/mozzy.git
cd mozzy
go build -o mozzy .
sudo mv mozzy /usr/local/bin/
```

---

## 🎯 Quick Start

```bash
# Simple GET request
mozzy GET https://jsonplaceholder.typicode.com/users/1

# Filter JSON with --jq
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city

# POST with JSON body
mozzy POST https://api.example.com/login --json '{"user":"alice","password":"secret"}'

# Verbose mode (see headers & timing)
mozzy GET https://api.example.com/users --verbose

# Save to collection
mozzy save my-api GET https://api.example.com/users --desc "Get all users"
mozzy exec my-api
```

---

## 📖 Features

### 🌐 HTTP Verbs

All the verbs you need with clean syntax:

```bash
mozzy GET /users
mozzy POST /users --json '{"name":"Alice"}'
mozzy PUT /users/1 --json '{"name":"Bob"}'
mozzy PATCH /users/1 --json '{"active":true}'
mozzy DELETE /users/1
```

### 🎨 Beautiful Output

Auto-colorized JSON that adapts to your terminal:
- **Cyan** keys
- **Green** strings
- **Yellow** numbers
- **Magenta** booleans
- **Auto TTY detection** - colors only when you need them

Force colors anywhere:
```bash
mozzy GET /api/data --color
```

### 🔍 JSONPath Filtering

Built-in JSON querying without external tools:

```bash
# Simple field
mozzy GET /users/1 --jq .name

# Nested path
mozzy GET /users/1 --jq .address.city

# Array indexing
mozzy GET /users --jq .[0].email

# Complex queries
mozzy GET /users --jq .[0].company.name
```

### 📚 Request Collections

Save and organize your API requests:

```bash
# Save requests
mozzy save github-user GET https://api.github.com/users/torvalds \
  --desc "Get Linus Torvalds profile"

# List collections
mozzy list

# Execute saved requests
mozzy exec github-user
```

### 🔗 API Chaining & Variables

Capture values from responses and use them in subsequent requests:

```bash
# Capture from response
mozzy GET /users --capture userId=[0].id

# Use captured variable
mozzy GET /users/{{userId}}/posts

# Chain authentication
mozzy POST /auth --json @creds.json --capture token=.access_token
mozzy GET /profile --auth "{{token}}"
```

### ⚙️ YAML Workflows

Automate multi-step API flows:

**workflow.yaml:**
```yaml
name: User Onboarding Flow
steps:
  - name: Create user
    method: POST
    url: https://api.example.com/users
    json: {"name": "Alice", "email": "alice@example.com"}
    capture:
      userId: .id

  - name: Send welcome email
    method: POST
    url: https://api.example.com/emails
    json:
      to: "alice@example.com"
      template: "welcome"
      userId: "{{userId}}"
```

```bash
mozzy run workflow.yaml
```

### 🔐 JWT Tools

Decode, verify, and sign JWTs without external tools:

```bash
# Decode JWT
mozzy jwt decode eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Verify with HMAC secret
mozzy jwt verify <token> --secret your-secret-key

# Verify with JWKS URL
mozzy jwt verify <token> --jwk https://example.com/.well-known/jwks.json

# Sign payload
echo '{"user":"alice","role":"admin"}' > payload.json
mozzy jwt sign payload.json --secret your-secret-key
```

### 🌍 Environment Management

Manage multiple environments (dev, staging, prod):

**.mozzy.json:**
```json
{
  "environments": {
    "dev": {
      "base_url": "http://localhost:3000",
      "headers": {"X-Env": "development"}
    },
    "staging": {
      "base_url": "https://staging.api.example.com",
      "auth_token": "staging-token"
    },
    "prod": {
      "base_url": "https://api.example.com",
      "auth_token": "prod-token"
    }
  }
}
```

```bash
# Use environment
mozzy --env prod GET /users

# List environments
mozzy env
```

### 📜 Request History

Browse and replay past requests:

```bash
# View recent requests
mozzy history

# Limit results
mozzy history --limit 10

# JSON output
mozzy history --json
```

### 🔧 Advanced Features

**Verbose Mode:**
```bash
mozzy GET /api/users --verbose
# Shows:
# - Request headers
# - Response headers
# - Timing breakdown (DNS, TLS, server, transfer)
```

**Retry Logic:**
```bash
# Retry on failure with exponential backoff
mozzy GET /flaky-endpoint --retry 3
```

**Cookie Jar:**
```bash
# Session persistence
mozzy GET /login --cookie-jar session.txt
mozzy GET /profile --cookie-jar session.txt
```

**Custom Headers:**
```bash
mozzy GET /api/data \
  --header "X-API-Key: secret123" \
  --header "X-Custom: value"
```

**Authentication:**
```bash
# Bearer token
mozzy GET /api/protected --auth your-token-here

# From environment
mozzy --env prod GET /api/protected
```

---

## 📚 Real-World Examples

### Example 1: GitHub API

```bash
# Get user info
mozzy GET https://api.github.com/users/torvalds --jq .name

# Get repositories
mozzy GET https://api.github.com/users/torvalds/repos --jq .[0].name

# Save for reuse
mozzy save gh-user GET https://api.github.com/users/torvalds
mozzy exec gh-user
```

### Example 2: Authentication Flow

```bash
# Login and capture token
mozzy POST https://api.example.com/auth \
  --json '{"username":"alice","password":"secret"}' \
  --capture token=.access_token

# Use token
mozzy GET https://api.example.com/profile \
  --auth "{{token}}" \
  --jq .email
```

### Example 3: CI/CD Testing

```bash
# Test API with --fail flag (exits non-zero on error)
mozzy GET https://api.example.com/health --fail

# Retry on failure
mozzy GET https://api.example.com/users --retry 3 --fail

# Use in scripts
if mozzy GET $API_URL/health --fail > /dev/null 2>&1; then
  echo "API is healthy"
else
  echo "API is down"
  exit 1
fi
```

### Example 4: Complex Workflow

**onboarding.yaml:**
```yaml
name: Complete User Onboarding
steps:
  - name: Register user
    method: POST
    url: https://api.example.com/register
    json:
      name: "Alice Smith"
      email: "alice@example.com"
    capture:
      userId: .user.id

  - name: Verify email
    method: POST
    url: https://api.example.com/verify
    json:
      userId: "{{userId}}"

  - name: Create profile
    method: POST
    url: https://api.example.com/profiles
    json:
      userId: "{{userId}}"
      bio: "Software Engineer"

  - name: Subscribe to newsletter
    method: POST
    url: https://api.example.com/subscriptions
    json:
      userId: "{{userId}}"
      type: "weekly"
```

```bash
mozzy run onboarding.yaml
```

---

## 🎨 Command Reference

### Global Flags

| Flag | Description |
|------|-------------|
| `--base <url>` | Base URL for requests |
| `--auth <token>` | Bearer authentication token |
| `--header <h:v>` | Custom header (repeatable) |
| `--env <name>` | Use named environment |
| `--jq <query>` | JSONPath filter |
| `--timeout <dur>` | Request timeout (default: 30s) |
| `--fail` | Exit non-zero on HTTP >= 400 |
| `--color` | Force colored output |
| `--no-color` | Disable colored output |
| `--verbose` / `-v` | Show headers & timing |
| `--retry <n>` | Retry attempts with backoff |
| `--cookie-jar <file>` | Cookie persistence file |
| `--capture <name=path>` | Capture variable (repeatable) |

### Commands

| Command | Description |
|---------|-------------|
| `GET/POST/PUT/DELETE/PATCH` | HTTP verbs |
| `save <name> <verb> <url>` | Save request to collection |
| `list` | List saved requests |
| `exec <name>` | Execute saved request |
| `history` | Show request history |
| `run <workflow.yaml>` | Run YAML workflow |
| `env` | List environments |
| `jwt decode <token>` | Decode JWT |
| `jwt verify <token>` | Verify JWT |
| `jwt sign <file>` | Sign JWT |

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

MIT © 2025 Archith Sharma

---

## 🌟 Show Your Support

If mozzy makes your API testing easier, give it a ⭐️ on GitHub!

**Found a bug?** [Open an issue](https://github.com/humancto/mozzy/issues)

**Have a feature request?** [Start a discussion](https://github.com/humancto/mozzy/discussions)

---

## 🔗 Links

- **Homepage:** https://github.com/humancto/mozzy
- **Releases:** https://github.com/humancto/mozzy/releases
- **Issues:** https://github.com/humancto/mozzy/issues
- **Discussions:** https://github.com/humancto/mozzy/discussions

---

<p align="center">
  <strong>Made with ❤️ by developers, for developers</strong>
</p>
