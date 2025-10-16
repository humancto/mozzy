# Changelog

All notable changes to mozzy will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.13.0] - 2025-10-16

### Added
- **HTTPS Proxy Support** (Phase 2) - Full SSL/TLS interception for HTTPS traffic
  - `--https` flag enables HTTPS interception mode
  - Automatic CA certificate generation and management in `~/.mozzy/`
  - Dynamic per-host certificate generation and caching
  - TLS handshake handling with man-in-the-middle proxy
  - `--export-cert` flag to export CA certificate for installation
  - `--cert-info` flag to view CA certificate details
  - Seamless integration with existing HTTP proxy
  - Verbose mode shows TLS handshake progress
  - Certificate valid for 10 years (CA) and 1 year (per-host)

### Technical Details
- Self-signed CA certificate generated using 2048-bit RSA keys
- Per-host certificates signed by CA and cached for performance
- CONNECT method handling for TLS tunnel establishment
- Proper certificate chain (host cert + CA cert) for client validation
- No external dependencies required

### Examples
```bash
mozzy proxy 8888 --https              # Start HTTPS proxy
mozzy proxy --export-cert > mozzy-ca.pem  # Export CA certificate
mozzy proxy --cert-info               # View certificate info

# Use with curl
curl -x http://localhost:8888 --cacert mozzy-ca.pem https://api.example.com

# Configure browser (after installing CA certificate)
HTTPS Proxy: localhost:8888
Port: 8888
```

### Installation Guide
To use HTTPS interception, you must install the CA certificate:
1. Export: `mozzy proxy --export-cert > mozzy-ca.pem`
2. macOS: Add to Keychain Access and mark as trusted
3. Linux: Add to `/usr/local/share/ca-certificates/` and run `update-ca-certificates`
4. Windows: Import into Trusted Root Certification Authorities

## [1.12.0] - 2025-10-16

### Added
- **HTTP Proxy Server** (Phase 1) - Intercept and inspect HTTP traffic
  - `mozzy proxy` command to start proxy server
  - Real-time traffic logging with timestamps
  - Colorized output for methods and status codes
  - Request duration tracking
  - Automatic local IP detection for easy mobile setup
  - Works with curl (`-x flag`), browsers, and any HTTP client
  - Request/response size tracking
  - Professional dashboard UI with box-drawing characters
  - Verbose mode for detailed request inspection

### Use Cases
- Debug API calls from any application
- Inspect third-party integrations
- Test mobile apps (HTTP only - HTTPS coming in Phase 2)
- Learn how HTTP requests work
- Capture traffic for analysis

### Examples
```bash
mozzy proxy 8888                    # Start on port 8888
mozzy proxy 8888 --verbose          # Detailed logging

# Use with curl
curl -x http://localhost:8888 http://api.example.com

# Configure browser
HTTP Proxy: localhost:8888
Port: 8888
```

### Coming Next
- Phase 2: HTTPS support with certificate generation
- Phase 3: Request/response modification, breakpoints, HAR recording

## [1.11.0] - 2025-10-16

### Added
- **Network Throttling** - Simulate different network speeds like Charles Proxy
  - 8 predefined profiles: `56k`, `slow`, `gprs`, `edge`, `3g`, `4g`, `lte`, `5g`
  - Realistic bandwidth limiting using token bucket algorithm
  - Latency simulation for each profile (10ms to 500ms)
  - `--throttle` flag available on all HTTP commands
  - Throttling info displayed in verbose mode
  - Perfect for testing slow network conditions, loading states, and mobile scenarios

### Examples
```bash
mozzy GET /api/large-file --throttle 3g      # 3G speed (~750 Kbps)
mozzy GET /api/data --throttle slow --verbose # Slow connection with details
mozzy POST /api/users --throttle 56k         # Dial-up speed testing
```

## [1.10.1] - 2025-10-16

### Changed
- **Enhanced Diff Visual Output** - Improved JSON diff display formatting
  - Box-drawing characters for professional headers
  - Color-coded symbols (+ for added, - for removed, ~ for changed, ! for type mismatch)
  - Better value formatting with proper quoting
  - Visual separators with Unicode characters
  - Git-style diff presentation

## [1.10.0] - 2025-10-16

### Added
- **JSON Response Diffing** - Compare API responses between environments
  - `mozzy diff` command for comparing two JSON files
  - Recursive comparison of nested objects and arrays
  - Color-coded differences (green for added, red for removed, yellow for changed)
  - Type mismatch detection (magenta)
  - Shows added fields, removed fields, changed values, and type changes
  - Perfect for detecting API contract changes between staging and production

### Examples
```bash
mozzy diff prod-response.json staging-response.json
curl https://prod.api.com/users/1 > prod.json
curl https://staging.api.com/users/1 > staging.json
mozzy diff prod.json staging.json
```

## [1.9.0] - 2025-10-15

### Added
- **Mock HTTP Server** - Built-in API mocking server
  - `mozzy mock` command to start HTTP mock server
  - YAML configuration support with routes, responses, headers
  - Generate sample config with `--generate` flag
  - Use saved requests as mocks with `--from-collection`
  - CORS support and custom headers
  - Response delays for testing timeouts
  - In-memory request logging
  - Perfect for frontend development and testing

### Examples
```bash
mozzy mock 8080 --generate > mock.yaml     # Generate sample config
mozzy mock 8080 --config mock.yaml         # Start with config
mozzy mock 3000 --from-collection          # Use saved requests as mocks
```

## [1.8.1] - 2025-10-15

### Added
- **Interactive Mode** - Browse history and saved requests with arrow keys
  - `mozzy interactive` or `mozzy i` command
  - `--saved` flag to browse saved collections
  - `--history` flag to browse request history
  - Arrow key navigation through requests
  - Press Enter to execute selected request
  - ESC or Ctrl+C to exit
  - Beautiful TUI interface with colored highlighting

### Examples
```bash
mozzy i                  # Browse history interactively
mozzy i --saved          # Browse saved collections
mozzy interactive        # Long form alias
```

## [1.8.0] - 2025-10-15

### Added
- **Performance Grading** - A-F letter grades for request performance
  - DNS Lookup grading (<50ms=A, 50-100ms=B, 100-200ms=C, >200ms=D/F)
  - TCP Connection grading (<30ms=A, 30-70ms=B, 70-150ms=C, >150ms=D/F)
  - TLS Handshake grading (<50ms=A, 50-120ms=B, 120-250ms=C, >250ms=D/F)
  - Server Response (TTFB) grading (<100ms=A, 100-300ms=B, 300-1000ms=C, >1000ms=D/F)
  - Overall performance grade calculation
  - Performance insights and optimization recommendations
  - Color-coded grades (green=A, blue=B, yellow=C, red=D/F)

### Enhanced
- **Verbose Mode** now includes:
  - Visual timeline with progress bars showing percentage of total time
  - Performance grade section with all timing grades
  - Intelligent optimization tips based on bottlenecks
  - Enhanced formatting and emojis

## [1.7.0] - 2025-10-14

### Added
- **Update Command** - Check for new mozzy versions
  - `mozzy update` command to check GitHub releases
  - Compares current version with latest release
  - Shows release notes and download link
  - Beautiful formatted output with version comparison

### Examples
```bash
mozzy update              # Check for updates
mozzy version             # Show current version
```

## [1.6.0] - 2025-10-14

### Added
- **Beautiful UI with Lipgloss** - Complete visual overhaul
  - Box-drawing characters for better visual hierarchy
  - Rounded borders and styled banners
  - Color-coded success, error, warning, and info messages
  - Professional table layout for collections list
  - Themed color palette (blue, green, red, yellow, cyan, purple)
- **New UI Package** (`internal/ui`)
  - Reusable styled components (banners, boxes, tables)
  - Success/Error/Warning/Info banner functions
  - Table rendering with box-drawing characters
  - Key-value pair formatting helpers
  - Help screen utilities with categorized commands
- **Enhanced Commands**
  - `mozzy list` - Beautiful table with borders showing Name, Method, URL, Description
  - `mozzy save` - Success banner with checkmark on save
  - Better empty state messages with helpful tips

### Changed
- Collections list now displays in a formatted table instead of plain text
- Success messages use styled banners with emoji and borders
- Improved visual consistency across all commands
- Added helpful tips after command outputs

### Documentation
- Added `VISUAL_IMPROVEMENTS.md` - Comprehensive UI enhancement brainstorm
- Added `DEMOS.md` - Guide for creating and updating demo GIFs
- Created 5 VHS demo recordings showing key features
- Updated README with showcase GIF and example demonstrations

## [1.5.0] - 2025-10-14

### Added
- **File Uploads** - `mozzy upload` command with multipart form support
  - Single and multiple file uploads
  - Custom field names for each file
  - Form data fields with `--data` flag
  - Progress tracking with spinner animation
  - Custom headers and bearer token authentication
  - Optional `--no-progress` for scripting
  - Human-readable file sizes and upload time
- Comprehensive test suite for upload functionality (40+ tests)
- Example scripts for various upload scenarios

## [1.4.0] - 2025-10-14

### Added
- **File Downloads** - `mozzy download` command with progress tracking
  - Real-time progress bar with ETA and download speed
  - Auto-detect filename from URL or Content-Disposition header
  - Overwrite protection with `--overwrite` flag
  - Support for large file streaming
  - Custom output path with `-o` flag
  - Optional progress display with `--no-progress`
  - Human-readable file sizes (KB, MB, GB)
  - Elapsed time and speed calculation
- Comprehensive test suite for download functionality (40+ tests)
- Example scripts for file download scenarios

## [1.3.0] - 2025-10-14

### Added
- **Conditional Retry** - `--retry-on` flag for fine-grained retry control
  - Retry on specific status codes (`--retry-on "503"`)
  - Retry on status ranges (`--retry-on "5xx"`, `--retry-on ">=500"`)
  - Retry on multiple conditions (`--retry-on "429,5xx,network_error"`)
  - Comparison operators: `>=`, `<=`, `>`, `<`, `==`, `!=`
  - Special conditions: `always`, `never`, `network_error`
- **Schema Validation** - Validate JSON responses against JSON Schema
  - Support for type validation (string, number, integer, boolean, array, object)
  - Object property validation with required fields
  - Nested object and array validation
  - String constraints (minLength, maxLength)
  - Number constraints (minimum, maximum)
  - Enum validation
  - Additional properties control
- **Conditional Workflows** - Control flow in YAML workflows
  - `on_success`: action to take on step success (continue, stop, or jump to step)
  - `on_failure`: action to take on step failure (stop, continue, or jump to step)
  - Enables retry loops, error handling, and complex branching logic
  - Default behavior: continue on success, stop on failure
- Comprehensive test suites for all new features (100+ tests)
- Example workflows and usage scripts

## [1.2.0] - 2025-10-14

### Added
- **Load Testing** - `mozzy load` command for performance testing
  - Fixed request count or duration-based testing
  - Configurable concurrency (`--concurrent`)
  - Detailed metrics: requests/sec, min/max/avg response times
  - Real-time progress reporting
- **Export Functionality** - `mozzy export` command to convert saved requests
  - Export to curl commands (`--format curl`)
  - Export to Postman collections (`--format postman`)
  - Works with both saved collections and workflows

## [1.1.0] - 2025-10-14

### Added
- **Test Suites** - `mozzy test` command to run workflows as automated tests
  - CI-friendly exit codes (0 for pass, 1 for fail)
  - JUnit XML output support (`--junit-output`)
  - Pass/fail summary with timing
- **Response Diffing** - `mozzy diff` command to compare JSON responses
  - Visual diff with color-coded changes
  - Deep JSON comparison
  - Useful for comparing environments
- Flow.Description field for better workflow documentation
- Examples for testing and diffing in `examples/testing/`

## [1.0.2] - 2025-10-14

### Added
- **Response Assertions** - Test APIs directly in workflows
  - Status code validation (`status == 200`, `status >= 200`)
  - Response time checks (`response_time < 500ms`)
  - JSON path assertions (`.name == "Alice"`)
  - String contains (`.email contains "@example.com"`)
  - Field existence (`.id exists`)
  - Length validation (`length(.items) > 0`)
  - Array access (`.items[0].id == 1`)
- Comprehensive test suite for assertions (100% coverage)
- Example workflows with assertions
- Comparison table in README (mozzy vs curl/httpie/Postman)

### Fixed
- Workflow variable substitution now works correctly (`{{vars}}` properly interpolate)
- JWT verify now shows success message and expiration info
- Improved JSON colorization compatibility across terminals

### Documentation
- Added `examples/` directory with workflow, JWT, and collection examples
- Added `examples/workflows/test-with-assertions.yaml`
- Updated README with feature comparison table

## [1.0.1] - 2025-10-13

### Fixed
- JSON colorization now displays properly in all terminals
- Fixed ANSI escape codes showing as raw text
- Changed from Sprint() to Print()/Printf() for proper color rendering

### Added
- Version command (`mozzy version`)
- Better error messages for color issues

## [1.0.0] - 2025-10-13

### Added
- Initial release
- HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Colored JSON output
- Inline JQ filtering (`--jq`)
- Request collections (save, list, exec)
- YAML workflows with variable capture
- JWT tools (decode, verify, sign)
- Request history
- Environment management
- Verbose mode with timing breakdown
- Cookie jar support
- Retry with exponential backoff

[1.13.0]: https://github.com/humancto/mozzy/compare/v1.12.0...v1.13.0
[1.12.0]: https://github.com/humancto/mozzy/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/humancto/mozzy/compare/v1.10.1...v1.11.0
[1.10.1]: https://github.com/humancto/mozzy/compare/v1.10.0...v1.10.1
[1.10.0]: https://github.com/humancto/mozzy/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/humancto/mozzy/compare/v1.8.1...v1.9.0
[1.8.1]: https://github.com/humancto/mozzy/compare/v1.8.0...v1.8.1
[1.8.0]: https://github.com/humancto/mozzy/compare/v1.7.0...v1.8.0
[1.7.0]: https://github.com/humancto/mozzy/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/humancto/mozzy/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/humancto/mozzy/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/humancto/mozzy/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/humancto/mozzy/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/humancto/mozzy/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/humancto/mozzy/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/humancto/mozzy/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/humancto/mozzy/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/humancto/mozzy/releases/tag/v1.0.0
