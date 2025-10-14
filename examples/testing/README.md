# Testing Examples

Examples showing how to use mozzy for API testing.

## Test Suites

Run workflows as automated test suites:

```bash
# Run tests
mozzy test api-test-suite.yaml

# Generate JUnit XML for CI/CD
mozzy test api-test-suite.yaml --junit-output results.xml
```

## Response Diffing

Compare API responses:

```bash
# Compare two JSON files
mozzy diff prod-response.json staging-response.json

# Compare responses from different environments
mozzy GET /api/users/1 --env prod > prod.json
mozzy GET /api/users/1 --env staging > staging.json
mozzy diff prod.json staging.json
```

## Files

- `api-test-suite.yaml` - Complete test suite example
- `compare-environments.sh` - Script to compare API responses across environments
