Output ../assets/demo-05-post.gif

Require mozzy

Set FontSize 16
Set Width 1200
Set Height 600
Set Padding 20
Set Theme "Molokai"

Type "# Example 5: POST Request with JSON" Sleep 500ms Enter
Sleep 1s

Type "mozzy POST https://jsonplaceholder.typicode.com/posts \\" Enter
Type "  --json '{\"title\":\"Hello\",\"body\":\"Testing mozzy\",\"userId\":1}' \\" Enter
Type "  --color"
Sleep 1s
Enter
Sleep 3s

Type "# ✅ Created new post!" Sleep 500ms Enter
Sleep 2s
