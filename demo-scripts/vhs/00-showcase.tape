Output ../assets/demo-showcase.gif

Require mozzy

Set FontSize 18
Set FontFamily "JetBrains Mono"
Set Width 1400
Set Height 900
Set Padding 30
Set Theme "Dracula"
Set TypingSpeed 100ms

Hide
Type "clear"
Enter
Show

# Banner
Type "echo ''" Enter
Type "echo ' 🦣 MOZZY - Postman for your Terminal'" Enter
Type "echo ' ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━'" Enter
Type "echo ''" Enter
Sleep 1s

# Show help
Type "mozzy --help"
Sleep 500ms
Enter
Sleep 4s

Hide
Type "clear"
Enter
Show

# Example 1: Simple GET with colors
Type "echo '📍 Example 1: Beautiful colored JSON'" Enter
Sleep 1s
Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --color"
Sleep 500ms
Enter
Sleep 3s

Hide
Type "clear"
Enter
Show

# Example 2: JQ filtering
Type "echo '🔍 Example 2: Inline JQ filtering'" Enter
Sleep 1s
Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name --color"
Sleep 500ms
Enter
Sleep 2s

Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city --color"
Sleep 500ms
Enter
Sleep 3s

Hide
Type "clear"
Enter
Show

# Example 3: Save to collection
Type "echo '📚 Example 3: Request Collections'" Enter
Sleep 1s
Type "mozzy save demo GET https://api.github.com/users/torvalds --desc 'GitHub user'"
Sleep 500ms
Enter
Sleep 2s

Type "mozzy list"
Sleep 500ms
Enter
Sleep 3s

Hide
Type "clear"
Enter
Show

# Example 4: Verbose mode
Type "echo '⏱️  Example 4: Verbose mode with timing'" Enter
Sleep 1s
Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --verbose --color"
Sleep 500ms
Enter
Sleep 5s

Hide
Type "clear"
Enter
Show

# Closing
Type "echo ''" Enter
Type "echo ' ✨ Learn more: github.com/humancto/mozzy'" Enter
Type "echo ' ⭐ Star the repo!'" Enter
Type "echo ''" Enter
Sleep 3s
