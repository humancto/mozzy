<p align="center">
  <img src="assets/logo.png" alt="Mozzy Logo" width="200"/>
</p>

<h1 align="center">mozzy — Postman for your Terminal 🚀</h1>

<p align="center">
  <a href="https://github.com/humancto/mozzy/releases"><img src="https://img.shields.io/github/v/release/humancto/mozzy" alt="GitHub release"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" alt="Go Version"></a>
</p>

A modern, developer-friendly HTTP client built for the terminal. Think **Postman meets curl** with beautiful colors, inline JSON queries, request collections, JWT tools, and powerful workflow automation.

<p align="center">
  <img src="assets/demo-showcase.gif" alt="Mozzy Demo" width="100%"/>
</p>

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

### mozzy vs The Rest

| Feature | curl | httpie | Postman | **mozzy** |
|---------|:----:|:------:|:-------:|:---------:|
| Colored JSON | ❌ | ✅ | ✅ | ✅ |
| Inline JQ Queries | ❌ | ❌ | ❌ | ✅ |
| Request Collections | ❌ | ❌ | ✅ | ✅ |
| YAML Workflows | ❌ | ❌ | ✅ | ✅ |
| API Chaining | ❌ | ❌ | ⚠️ | ✅ |
| JWT Tools Built-in | ❌ | ❌ | ❌ | ✅ |
| Request History | ❌ | ❌ | ✅ | ✅ |
| CLI First | ✅ | ✅ | ❌ | ✅ |
| Free & Open Source | ✅ | ✅ | 💰 | ✅ |
| Easy to Learn | ❌ | ✅ | ✅ | ✅ |

---

## 🚀 Installation

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/humancto/mozzy/main/install.sh | bash
```

### Homebrew

```bash
# Tap the repository
brew tap humancto/mozzy

# Install (use full path to avoid typos)
brew install humancto/mozzy/mozzy
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

Automate multi-step API flows with conditional execution:

**workflow.yaml:**
```yaml
name: User Onboarding Flow with Error Handling
steps:
  - name: Create user
    method: POST
    url: https://api.example.com/users
    json: {"name": "Alice", "email": "alice@example.com"}
    capture:
      userId: .id
    on_success: send_email
    on_failure: error_handler

  - name: send_email
    method: POST
    url: https://api.example.com/emails
    json:
      to: "alice@example.com"
      template: "welcome"
      userId: "{{userId}}"
    assert:
      - status == 200
    on_success: continue
    on_failure: error_handler

  - name: error_handler
    method: POST
    url: https://api.example.com/errors
    json: {"message": "Onboarding failed"}
```

```bash
mozzy run workflow.yaml
```

**Conditional Execution:**
- `on_success`: continue (default), stop, or jump to step name
- `on_failure`: stop (default), continue, or jump to step name

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

**Enhanced Verbose Mode (v1.7.0):**
```bash
mozzy GET /api/users --verbose
# Shows:
# - 🔍 DNS Resolution (resolved IP, latency)
# - 🌐 Connection Protocol (HTTP/1.1, HTTP/2, HTTP/3)
# - 🔐 TLS Certificate Info (subject, issuer, expiry)
# - 📦 Transfer Sizes (request/response with compression ratio)
# - ⏱️ Visual Timeline with progress bars and percentages
# - Request & Response headers
```

**Retry Logic with Conditions:**
```bash
# Retry on failure with exponential backoff (default: 5xx errors)
mozzy GET /flaky-endpoint --retry 3

# Retry on specific conditions
mozzy GET /api/data --retry 5 --retry-on "429,5xx"  # Rate limit + server errors
mozzy GET /api/data --retry 3 --retry-on ">=500"    # Any status >= 500
mozzy GET /api/data --retry 3 --retry-on "503"      # Only 503
mozzy GET /api/data --retry 3 --retry-on "network_error"  # Network failures only
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

**File Downloads:**
```bash
# Download with progress bar
mozzy download https://example.com/file.zip

# Custom output path
mozzy download https://example.com/file.zip -o myfile.zip

# Overwrite existing files
mozzy download https://example.com/file.zip --overwrite

# Disable progress for scripting
mozzy download https://example.com/data.json --no-progress
```

**File Uploads:**
```bash
# Upload single file
mozzy upload https://api.example.com/upload -f avatar.jpg

# Upload multiple files
mozzy upload https://api.example.com/upload -f file1.jpg -f file2.png

# Upload with form data
mozzy upload https://api.example.com/upload -f resume.pdf --data "name=John" --data "email=john@example.com"

# Upload with authentication
mozzy upload https://api.example.com/upload -f file.jpg --auth token123
```

---

## 📚 Real-World Examples

### 🐙 Example 1: Exploring GitHub's API

**Use case**: Query GitHub's public API to fetch user profiles and repositories.

<p align="center">
  <img src="assets/demo-01-github-api.gif" alt="GitHub API Demo" width="90%"/>
</p>

```bash
# Get user info with colored output
mozzy GET https://api.github.com/users/torvalds --color

# Extract just the name
mozzy GET https://api.github.com/users/torvalds --jq .name

# Get first repository name
mozzy GET https://api.github.com/users/torvalds/repos --jq .[0].name

# Save for reuse in your collection
mozzy save gh-torvalds GET https://api.github.com/users/torvalds \
  --desc "Get Linus Torvalds GitHub profile"

# Execute anytime
mozzy exec gh-torvalds
```

---

### 🔐 Example 2: OAuth & Authentication Workflows

**Use case**: Login to an API, capture the access token, and use it for authenticated requests.

```bash
# Step 1: Login and capture the token
mozzy POST https://api.example.com/auth/login \
  --json '{"username":"alice","password":"secret123"}' \
  --capture token=.access_token \
  --verbose

# Step 2: Use the captured token for authenticated requests
mozzy GET https://api.example.com/user/profile \
  --auth "{{token}}" \
  --jq .email

# Step 3: Update profile with the token
mozzy PATCH https://api.example.com/user/profile \
  --auth "{{token}}" \
  --json '{"bio":"DevOps Engineer"}' \
  --verbose
```

---

### ⚙️ Example 3: CI/CD Health Checks & Monitoring

**Use case**: Integrate mozzy into your deployment pipeline to verify API health and run smoke tests.

```bash
# Basic health check with exit code
mozzy GET https://api.example.com/health --fail

# Health check with retry logic (5xx errors)
mozzy GET https://api.example.com/health \
  --retry 3 \
  --retry-on "5xx" \
  --fail

# Advanced: Health check script for CI/CD
#!/bin/bash
if mozzy GET $API_URL/health --fail --timeout 10s > /dev/null 2>&1; then
  echo "✅ API is healthy"
  mozzy GET $API_URL/metrics --jq .uptime
else
  echo "❌ API is down - deploying rollback"
  exit 1
fi

# Load testing before production deployment
mozzy load https://staging.api.example.com/users \
  --requests 1000 \
  --concurrent 10
```

---

### 🔗 Example 4: Multi-Step E-commerce Workflow

**Use case**: Automate a complete user onboarding flow with conditional error handling.

**onboarding-flow.yaml:**
```yaml
name: E-commerce User Onboarding
description: Register user, verify email, create profile, and subscribe

steps:
  - name: register
    method: POST
    url: https://api.shop.example.com/register
    json:
      name: "Alice Smith"
      email: "alice@example.com"
      password: "secure123"
    capture:
      userId: .user.id
      verifyToken: .verification_token
    assert:
      - status == 201
      - .user.id exists
    on_success: verify_email
    on_failure: error_handler

  - name: verify_email
    method: POST
    url: https://api.shop.example.com/verify
    json:
      userId: "{{userId}}"
      token: "{{verifyToken}}"
    assert:
      - status == 200
      - .verified == true
    on_success: create_profile
    on_failure: error_handler

  - name: create_profile
    method: POST
    url: https://api.shop.example.com/profiles
    json:
      userId: "{{userId}}"
      bio: "Software Engineer"
      preferences:
        newsletter: true
        notifications: true
    capture:
      profileId: .profile.id
    on_success: subscribe

  - name: subscribe
    method: POST
    url: https://api.shop.example.com/subscriptions
    json:
      userId: "{{userId}}"
      plan: "premium"
      billing: "monthly"
    assert:
      - status == 200
      - .subscription.active == true

  - name: error_handler
    method: POST
    url: https://api.shop.example.com/errors/log
    json:
      flow: "onboarding"
      message: "User onboarding failed"
      userId: "{{userId}}"
```

```bash
# Run the complete workflow
mozzy run onboarding-flow.yaml

# Run as test suite with JUnit output for CI
mozzy test onboarding-flow.yaml --junit-output results.xml
```

---

### 📦 Example 5: File Upload & Download Operations

**Use case**: Upload product images to an e-commerce API and download invoices.

```bash
# Upload single product image
mozzy upload https://api.shop.example.com/products/123/images \
  -f product-photo.jpg \
  --auth "your-token" \
  --data "alt_text=Blue Widget" \
  --data "primary=true"

# Upload multiple images at once
mozzy upload https://api.shop.example.com/products/456/gallery \
  -f image1.jpg \
  -f image2.jpg \
  -f image3.png \
  --auth "your-token"

# Download invoice with progress bar
mozzy download https://api.shop.example.com/invoices/789/pdf \
  -o invoice-789.pdf \
  --auth "your-token"

# Download large dataset without progress (for scripts)
mozzy download https://api.example.com/exports/data.zip \
  --no-progress \
  --overwrite
```

---

### 🧪 Example 6: API Testing & Validation

**Use case**: Run automated API tests with assertions and schema validation.

```bash
# Test with inline assertions
mozzy GET https://api.example.com/users/1 \
  --assert "status == 200" \
  --assert ".name exists" \
  --assert ".email contains @example.com"

# Compare responses between environments
mozzy GET https://staging.api.example.com/users/1 -o staging.json
mozzy GET https://prod.api.example.com/users/1 -o prod.json
mozzy diff staging.json prod.json

# Export collection to Postman for team sharing
mozzy export my-api-collection --format postman -o postman-collection.json

# Export to curl for documentation
mozzy export my-request --format curl
```

**test-suite.yaml:**
```yaml
name: API Test Suite
steps:
  - name: Test user creation
    method: POST
    url: https://api.example.com/users
    json:
      name: "Test User"
      email: "test@example.com"
    assert:
      - status == 201
      - .id exists
      - .email == "test@example.com"
      - response_time < 500ms

  - name: Test user retrieval
    method: GET
    url: https://api.example.com/users/1
    assert:
      - status == 200
      - .name exists
      - length(.friends) >= 0
```

```bash
# Run test suite
mozzy test test-suite.yaml --junit-output test-results.xml
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
| `--retry-on <cond>` | Retry conditions (5xx, 429, >=500, etc.) |
| `--cookie-jar <file>` | Cookie persistence file |
| `--capture <name=path>` | Capture variable (repeatable) |

### Commands

| Command | Description |
|---------|-------------|
| `GET/POST/PUT/DELETE/PATCH` | HTTP verbs |
| `download <url>` | Download files with progress bar |
| `upload <url>` | Upload files with multipart forms |
| `save <name> <verb> <url>` | Save request to collection |
| `list` | List saved requests |
| `exec <name>` | Execute saved request |
| `history` | Show request history |
| `run <workflow.yaml>` | Run YAML workflow |
| `test <workflow.yaml>` | Run workflow as test suite |
| `diff <file1> <file2>` | Compare JSON responses |
| `load <url>` | Performance load testing |
| `export <name>` | Export to curl/Postman |
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

MIT © 2025 Archy

---

## 🌟 Show Your Support

If mozzy makes your API testing easier, give it a ⭐️ on GitHub!

**Found a bug?** [Open an issue](https://github.com/humancto/mozzy/issues)

**Have a feature request?** [Start a discussion](https://github.com/humancto/mozzy/discussions)

**Enjoying mozzy?** [Buy me a coffee](https://buymeacoffee.com/humancto) ☕

---

## 👨‍💻 About the Author

Built by **Archy** ([@humancto](https://github.com/humancto))

- 🌐 Website: [humancto.com](https://www.humancto.com)
- ☕ Support: [buymeacoffee.com/humancto](https://buymeacoffee.com/humancto)
- 💼 GitHub: [@humancto](https://github.com/humancto)

---

## 🔗 Links

- **Homepage:** https://github.com/humancto/mozzy
- **Releases:** https://github.com/humancto/mozzy/releases
- **Issues:** https://github.com/humancto/mozzy/issues
- **Discussions:** https://github.com/humancto/mozzy/discussions

---

<p align="center">
  <strong>Made with ❤️ by <a href="https://www.humancto.com">humancto.com</a>, for developers</strong>
</p>
