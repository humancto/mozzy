# Changelog

All notable changes to mozzy will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[1.4.0]: https://github.com/humancto/mozzy/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/humancto/mozzy/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/humancto/mozzy/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/humancto/mozzy/compare/v1.0.2...v1.1.0
[1.0.2]: https://github.com/humancto/mozzy/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/humancto/mozzy/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/humancto/mozzy/releases/tag/v1.0.0
