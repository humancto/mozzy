package formatter

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type StatusInfo struct {
	Emoji       string
	Description string
	Tips        []string
}

var statusMessages = map[int]StatusInfo{
	// 2xx Success
	200: {Emoji: "✅", Description: "OK - Request successful", Tips: nil},
	201: {Emoji: "✨", Description: "Created - Resource successfully created", Tips: nil},
	202: {Emoji: "⏳", Description: "Accepted - Request accepted for processing", Tips: nil},
	204: {Emoji: "📭", Description: "No Content - Successful but no content to return", Tips: nil},

	// 3xx Redirection
	301: {Emoji: "↪️", Description: "Moved Permanently", Tips: []string{"Update your bookmark/URL to the new location"}},
	302: {Emoji: "🔀", Description: "Found (Temporary Redirect)", Tips: nil},
	304: {Emoji: "💾", Description: "Not Modified - Cached version is still valid", Tips: nil},
	307: {Emoji: "🔄", Description: "Temporary Redirect", Tips: nil},
	308: {Emoji: "⏩", Description: "Permanent Redirect", Tips: nil},

	// 4xx Client Errors
	400: {Emoji: "❌", Description: "Bad Request - Invalid syntax or malformed request", Tips: []string{
		"Check your JSON payload for syntax errors",
		"Verify all required fields are present",
		"Ensure data types match API expectations",
	}},
	401: {Emoji: "🔒", Description: "Unauthorized - Authentication required or failed", Tips: []string{
		"Check if you've provided an auth token: --auth <token>",
		"Verify your token hasn't expired",
		"Ensure you're using the correct authentication method (Bearer, API Key, etc.)",
		"Try regenerating your access token",
	}},
	403: {Emoji: "🚫", Description: "Forbidden - Authenticated but not authorized", Tips: []string{
		"You're logged in but don't have permission to access this resource",
		"Check if your account has the required role/permissions",
		"Verify you're accessing the correct resource ID",
		"Contact your admin to grant necessary permissions",
	}},
	404: {Emoji: "🔍", Description: "Not Found - Resource doesn't exist", Tips: []string{
		"Double-check the URL/endpoint path",
		"Verify the resource ID is correct",
		"Check if the resource was deleted",
		"Ensure you're using the right base URL (--base or --env)",
	}},
	405: {Emoji: "🛑", Description: "Method Not Allowed", Tips: []string{
		"This endpoint doesn't support the HTTP method you used",
		"Try a different method (GET instead of POST, etc.)",
		"Check the API documentation for allowed methods",
	}},
	408: {Emoji: "⏱️", Description: "Request Timeout", Tips: []string{
		"The server took too long to respond",
		"Try increasing timeout: --timeout 60s",
		"Check your network connection",
	}},
	409: {Emoji: "⚠️", Description: "Conflict - Request conflicts with current state", Tips: []string{
		"Resource already exists or has been modified",
		"Check if you're trying to create a duplicate",
		"Refresh and try again with updated data",
	}},
	415: {Emoji: "📄", Description: "Unsupported Media Type", Tips: []string{
		"Server doesn't support the content type you sent",
		"Add proper Content-Type header: --header 'Content-Type: application/json'",
		"Ensure you're sending data in the expected format",
	}},
	422: {Emoji: "🔧", Description: "Unprocessable Entity - Validation failed", Tips: []string{
		"Request syntax is correct but semantically invalid",
		"Check field validation rules (min/max length, format, etc.)",
		"Review error details in response body for specific field errors",
	}},
	429: {Emoji: "🐌", Description: "Too Many Requests - Rate limit exceeded", Tips: []string{
		"You've hit the API rate limit",
		"Wait a few minutes before retrying",
		"Check response headers for rate limit reset time",
		"Consider implementing exponential backoff",
	}},

	// 5xx Server Errors
	500: {Emoji: "💥", Description: "Internal Server Error - Something went wrong on the server", Tips: []string{
		"This is a server-side issue, not your fault",
		"Try again in a few moments",
		"Check API status page if available",
		"Contact API support if error persists",
	}},
	501: {Emoji: "🚧", Description: "Not Implemented - Feature not supported yet", Tips: []string{
		"This feature isn't implemented on the server",
		"Check API documentation for supported features",
		"Look for alternative endpoints",
	}},
	502: {Emoji: "🌐", Description: "Bad Gateway - Invalid response from upstream server", Tips: []string{
		"Gateway/proxy received invalid response",
		"This is usually temporary, try again",
		"Check if there's ongoing maintenance",
	}},
	503: {Emoji: "🔧", Description: "Service Unavailable - Server temporarily down", Tips: []string{
		"Server is temporarily unavailable (maintenance or overloaded)",
		"Check 'Retry-After' header for when to retry",
		"Visit API status page for updates",
		"Try again in a few minutes",
	}},
	504: {Emoji: "⏰", Description: "Gateway Timeout - Upstream server didn't respond in time", Tips: []string{
		"Gateway/proxy timeout waiting for server",
		"Try again with a longer timeout: --timeout 120s",
		"This might indicate server performance issues",
	}},
}

func PrintStatusLine(method, url string, statusCode int, duration time.Duration) {
	// Color codes based on status
	var statusColor *color.Color
	var arrow string

	switch {
	case statusCode >= 200 && statusCode < 300:
		statusColor = color.New(color.FgGreen, color.Bold)
		arrow = color.GreenString("→")
	case statusCode >= 300 && statusCode < 400:
		statusColor = color.New(color.FgCyan, color.Bold)
		arrow = color.CyanString("↪")
	case statusCode >= 400 && statusCode < 500:
		statusColor = color.New(color.FgYellow, color.Bold)
		arrow = color.YellowString("⚠")
	case statusCode >= 500:
		statusColor = color.New(color.FgRed, color.Bold)
		arrow = color.RedString("✗")
	default:
		statusColor = color.New(color.FgWhite)
		arrow = "→"
	}

	// Method color
	methodColor := color.New(color.FgMagenta, color.Bold)

	// Duration color (green if fast, yellow if slow, red if very slow)
	var durationColor *color.Color
	ms := duration.Milliseconds()
	switch {
	case ms < 200:
		durationColor = color.New(color.FgGreen)
	case ms < 1000:
		durationColor = color.New(color.FgYellow)
	default:
		durationColor = color.New(color.FgRed)
	}

	fmt.Fprintf(color.Output, "%s %s %s %s %s\n",
		arrow,
		methodColor.Sprint(method),
		url,
		statusColor.Sprintf("(%d)", statusCode),
		durationColor.Sprintf("in %s", duration),
	)

	// Print status explanation for non-2xx responses
	if statusCode < 200 || statusCode >= 300 {
		PrintStatusExplanation(statusCode)
	}
}

func PrintStatusExplanation(statusCode int) {
	info, exists := statusMessages[statusCode]
	if !exists {
		// Generic message for unknown status codes
		if statusCode >= 400 && statusCode < 500 {
			info = StatusInfo{
				Emoji:       "❓",
				Description: "Client Error",
				Tips:        []string{"Check the API documentation for details about this status code"},
			}
		} else if statusCode >= 500 {
			info = StatusInfo{
				Emoji:       "❓",
				Description: "Server Error",
				Tips:        []string{"This is a server-side issue", "Try again later or contact API support"},
			}
		} else {
			return
		}
	}

	// Print emoji and description
	titleColor := color.New(color.FgCyan, color.Bold)
	fmt.Fprintf(color.Output, "\n%s  %s\n", info.Emoji, titleColor.Sprint(info.Description))

	// Print tips if available
	if len(info.Tips) > 0 {
		tipColor := color.New(color.FgYellow)
		fmt.Fprintf(color.Output, "\n%s\n", color.New(color.FgWhite, color.Bold).Sprint("💡 Tips:"))
		for _, tip := range info.Tips {
			fmt.Fprintf(color.Output, "  • %s\n", tipColor.Sprint(tip))
		}
	}
	fmt.Println()
}
