# Verbose Output Enhancements

Ideas for making `--verbose` mode more useful for developers and debugging.

## 🌐 Network & Connection Info

### DNS Resolution
```
🔍 DNS Resolution
  • Domain:    api.example.com
  • IP:        93.184.216.34
  • Type:      IPv4
  • Resolver:  8.8.8.8 (Google DNS)
  • Latency:   12ms
```

### Connection Details
```
🔌 Connection
  • Protocol:  HTTPS/2.0
  • TLS:       TLS 1.3
  • Cipher:    TLS_AES_128_GCM_SHA256
  • Server:    nginx/1.21.0
  • Port:      443
  • Keep-Alive: enabled
```

### Certificate Info (for HTTPS)
```
🔐 TLS Certificate
  • Subject:   *.example.com
  • Issuer:    Let's Encrypt Authority X3
  • Valid:     2024-01-01 → 2025-01-01 (345 days remaining)
  • Algorithm: RSA 2048-bit
  • SNI:       api.example.com
```

## ⏱️ Enhanced Timing Breakdown

### Current (Good)
```
⏱ Timing Breakdown:
  • DNS Lookup:        2ms
  • TCP Connection:    24ms
  • TLS Handshake:     30ms
  • Server Processing: 83ms
  • Content Transfer:  22µs
  • Total:             83ms
```

### Enhanced (Better)
```
⏱ Request Timeline
  ┌─ DNS Lookup          2ms    ████░░░░░░░░░░░░░░░░
  ├─ TCP Connect         24ms   ████████████░░░░░░░░
  ├─ TLS Handshake       30ms   ███████████████░░░░░
  ├─ Request Sent        1ms    ░░░░░░░░░░░░░░░░░░░░
  ├─ Waiting (TTFB)      26ms   █████████░░░░░░░░░░░
  ├─ Content Download    22µs   ░░░░░░░░░░░░░░░░░░░░
  └─ Total              83ms    ████████████████████

  TTFB (Time to First Byte): 26ms
  Download Speed: 15.2 MB/s
```

## 📊 Request/Response Size Info

```
📦 Transfer Details
  Request:
    • Headers:   245 bytes
    • Body:      1.2 KB
    • Total:     1.4 KB

  Response:
    • Headers:   512 bytes
    • Body:      4.8 KB (gzip: 1.2 KB, 75% compression)
    • Total:     5.3 KB

  Bandwidth: 64 KB/s
```

## 🔍 HTTP/2 or HTTP/3 Specific Info

```
🚀 HTTP/2 Details
  • Stream ID:        1
  • Priority:         Weight 16
  • Server Push:      disabled
  • Multiplexing:     enabled
  • Window Size:      65535
```

## 🍪 Cookie Information

```
🍪 Cookies
  Request Cookies (2):
    • session_id: abc123... (HttpOnly, Secure)
    • preferences: theme=dark

  Response Cookies (1):
    • csrf_token: xyz789... (SameSite=Strict, expires in 1h)
```

## 🔄 Redirect Chain

```
🔄 Redirect Chain (3 hops)
  1. http://example.com        → 301 Moved Permanently (12ms)
  2. https://example.com       → 302 Found (45ms)
  3. https://www.example.com   → 200 OK (83ms)

  Total redirect time: 140ms
```

## 🛡️ Security Headers

```
🛡️ Security Analysis
  ✅ Strict-Transport-Security (HSTS enabled)
  ✅ Content-Security-Policy (XSS protection)
  ✅ X-Frame-Options (Clickjacking protection)
  ✅ X-Content-Type-Options (MIME sniffing disabled)
  ⚠️  Missing: Referrer-Policy
  ❌ Missing: Permissions-Policy

  Security Score: 4/6 (Good)
```

## 🌍 Geolocation Info

```
🌍 Server Location
  • IP:        93.184.216.34
  • Country:   United States 🇺🇸
  • Region:    Virginia
  • City:      Ashburn
  • ISP:       Amazon AWS
  • ASN:       AS16509
  • Distance:  ~2,800 km
```

## 🔁 Retry Information (when retries happen)

```
🔁 Retry Attempts
  Attempt 1: 503 Service Unavailable (failed) - waited 1s
  Attempt 2: 503 Service Unavailable (failed) - waited 2s
  Attempt 3: 200 OK (success)

  Total retry time: 3.2s
  Retries used: 2/3
```

## 🧬 Request Fingerprint

```
🧬 Request Fingerprint
  • Method:      GET
  • URL Hash:    a7f3c9d2...
  • Headers:     8 custom headers
  • User-Agent:  mozzy/1.6.0
  • Accept:      application/json
  • Compression: gzip, deflate, br
```

## 📈 Cache Information

```
💾 Cache Status
  • Status:      HIT
  • Age:         3695s (1h 1m)
  • Max-Age:     43200s (12h)
  • Remaining:   39505s (11h)
  • ETag:        W/"1fd-+2Y3G3w049..."
  • Validation:  If-None-Match sent
```

## 🎯 Performance Metrics

```
📊 Performance Metrics
  • DNS:            Fast (2ms < 20ms)
  • Connection:     Good (24ms < 50ms)
  • TLS:            Good (30ms < 100ms)
  • TTFB:           Excellent (26ms < 100ms)
  • Download:       Excellent (15MB/s)

  Overall Grade: A (Excellent)
```

## 🐛 Debug Information

```
🐛 Debug Info
  • Request ID:     req_abc123xyz
  • Trace ID:       trace_789def
  • Request UUID:   550e8400-e29b-41d4-a716-446655440000
  • Client IP:      192.168.1.100 (as seen by server)
  • Proxy:          None
```

## 🔗 Related Requests (for API chaining)

```
🔗 Variable Captures
  ✓ Captured 'userId' = 12345 from .id
  ✓ Captured 'token' = eyJhbGc... from .access_token

  Available for next request: {{userId}}, {{token}}
```

## 📋 Implementation Priority

### High Priority (Next Release - v1.7.0)
1. ✅ **DNS Resolution** - Show resolved IP address
2. ✅ **Connection Protocol** - HTTP/1.1, HTTP/2, HTTP/3
3. ✅ **TLS Version** - Security info
4. ✅ **Enhanced Timing** - Bar chart visualization
5. ✅ **Transfer Sizes** - Request/response sizes with compression info

### Medium Priority (v1.8.0)
6. **Certificate Info** - For HTTPS endpoints
7. **Redirect Chain** - Show full redirect path
8. **Security Headers** - Analyze and score
9. **Cache Information** - Cache hit/miss, age, etc.
10. **Performance Grading** - A-F grade based on metrics

### Low Priority (v2.0.0)
11. **Geolocation** - Server location (requires external API)
12. **HTTP/2 Details** - Stream info, server push, etc.
13. **Cookie Details** - Detailed cookie parsing
14. **Retry Visualization** - Show retry attempts with timeline

## 🎨 Visual Design

### Use Box Drawing for Sections
```
╭─ DNS Resolution ─────────────────────────────────╮
│ Domain:  api.example.com                         │
│ IP:      93.184.216.34 (IPv4)                    │
│ Latency: 12ms                                    │
╰──────────────────────────────────────────────────╯
```

### Use Progress Bars for Timing
```
DNS:        ████░░░░░░░░░░░░░░░░  12ms  (10%)
Connect:    ██████████░░░░░░░░░░  45ms  (37%)
TLS:        ████████░░░░░░░░░░░░  32ms  (26%)
TTFB:       ███████░░░░░░░░░░░░░  28ms  (23%)
Download:   █░░░░░░░░░░░░░░░░░░░  5ms   (4%)
```

### Use Colors for Status
```
✅ DNS Resolution:    2ms     (excellent)
✅ TCP Connection:    24ms    (good)
⚠️  TLS Handshake:    150ms   (slow)
✅ Server Response:   26ms    (excellent)
✅ Total:             202ms   (good)
```

## 💡 Implementation Notes

### Go Libraries to Use
```go
// DNS resolution
"net"
"net/netip"

// TLS info
"crypto/tls"
"crypto/x509"

// HTTP/2 details
"golang.org/x/net/http2"

// Timing
"time"
"github.com/charmbracelet/lipgloss" // For styling
```

### Store Connection Info
Capture these during request:
- `resp.TLS` - TLS connection state
- `resp.Request.RemoteAddr` - Server IP
- Custom RoundTripper to capture DNS, connection times
- `httptrace` package for detailed timing

### Example Code Structure
```go
type VerboseInfo struct {
    DNS           DNSInfo
    Connection    ConnectionInfo
    TLS           *tls.ConnectionState
    Timing        TimingInfo
    Sizes         SizeInfo
    Security      SecurityInfo
    Performance   PerformanceGrade
}

func RenderVerboseOutput(info VerboseInfo) string {
    var sb strings.Builder
    sb.WriteString(RenderDNSInfo(info.DNS))
    sb.WriteString(RenderConnectionInfo(info.Connection))
    sb.WriteString(RenderTimingChart(info.Timing))
    sb.WriteString(RenderSecurityScore(info.Security))
    return sb.String()
}
```

## 🎯 User Benefits

### For Developers
- Debug DNS issues quickly
- Identify slow connections vs slow servers
- Verify TLS/SSL configuration
- Optimize API performance

### For DevOps
- Monitor API latency from different locations
- Identify caching issues
- Check security header compliance
- Track redirect chains

### For Security Engineers
- Verify TLS versions and ciphers
- Check certificate validity
- Audit security headers
- Identify potential vulnerabilities

## 🚀 Future: Interactive Mode

```bash
mozzy GET /api --interactive
```

Opens a TUI (Terminal UI) with:
- Real-time request/response view
- Collapsible sections for headers, body, timing
- Copy-to-clipboard buttons
- Syntax highlighting
- Request history browser

Using `github.com/charmbracelet/bubbletea`
