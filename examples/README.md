# mozzy Examples

This directory contains practical examples demonstrating mozzy's features.

## Directory Structure

- **basic/** - Simple HTTP requests (GET, POST, PUT, DELETE, PATCH)
- **workflows/** - Multi-step API workflows with variable capture
- **jwt/** - JWT encoding, decoding, and verification examples
- **collections/** - Saved request collections

## Quick Start

```bash
# Run a basic GET request
mozzy GET https://jsonplaceholder.typicode.com/users/1

# Run a workflow
mozzy run examples/workflows/user-onboarding.yaml

# Decode a JWT
mozzy jwt decode $(cat examples/jwt/sample-token.txt)

# Execute a saved collection
cd examples/collections && mozzy exec get-user
```

## Example Categories

### 1. Basic HTTP Requests

Simple examples showing each HTTP method with various options.

### 2. Workflows

Multi-step API automation with variable capture and chaining.

### 3. JWT Tools

Examples of encoding, decoding, signing, and verifying JWTs.

### 4. Collections

Pre-configured request collections for common APIs.
