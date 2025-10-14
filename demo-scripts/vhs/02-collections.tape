Output ../assets/demo-02-collections.gif

Require mozzy

Set FontSize 16
Set Width 1200
Set Height 600
Set Padding 20
Set Theme "Molokai"

Type "# Example 2: Request Collections" Sleep 500ms Enter
Sleep 1s

Type "mozzy save gh-user GET https://api.github.com/users/torvalds --desc 'Get Linus profile'"
Sleep 1s
Enter
Sleep 2s

Type "mozzy list"
Sleep 1s
Enter
Sleep 3s

Type "mozzy exec gh-user"
Sleep 1s
Enter
Sleep 3s

Type "# ✅ Saved and executed!" Sleep 500ms Enter
Sleep 2s
