Output ../assets/demo-04-verbose.gif

Require mozzy

Set FontSize 14
Set Width 1200
Set Height 800
Set Padding 20
Set Theme "Molokai"

Type "# Example 4: Verbose Mode with Timing" Sleep 500ms Enter
Sleep 1s

Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --verbose --color"
Sleep 1s
Enter
Sleep 5s

Type "# ✅ See headers and timing breakdown!" Sleep 500ms Enter
Sleep 2s
