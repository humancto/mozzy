Output ../assets/demo-04-verbose.gif

Require mozzy

Set FontSize 13
Set Width 1400
Set Height 1000
Set Padding 20
Set Theme "Dracula"

Type "# Example 4: Enhanced Verbose Mode (v1.7.0)" Sleep 500ms Enter
Sleep 1s

Type "# 🔍 DNS Resolution, TLS Info, Transfer Sizes & Visual Timing!" Sleep 500ms Enter
Sleep 1.5s

Type "mozzy GET https://api.github.com/users/octocat --verbose --color"
Sleep 1s
Enter
Sleep 6s

Type "# ✅ See: IP, Protocol, TLS Cert, Sizes, Timeline bars!" Sleep 500ms Enter
Sleep 2s
