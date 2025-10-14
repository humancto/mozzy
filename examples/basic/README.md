# Basic Examples

Simple examples demonstrating core HTTP methods and features.

## Files

- `01-simple-get.sh` - GET requests with and without JQ filtering
- `02-post-request.sh` - POST requests with JSON bodies and headers
- `03-verbose-mode.sh` - Using verbose mode to see headers and timing

## Usage

```bash
# Make scripts executable
chmod +x examples/basic/*.sh

# Run any example
./examples/basic/01-simple-get.sh
```

## Examples

### Simple GET
```bash
mozzy GET https://jsonplaceholder.typicode.com/users/1
```

### GET with JQ Filter
```bash
mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name
```

### POST with JSON
```bash
mozzy POST https://jsonplaceholder.typicode.com/posts \
  --json '{"title":"Test","body":"Content","userId":1}'
```

### Verbose Mode
```bash
mozzy GET https://api.example.com/data --verbose
```
