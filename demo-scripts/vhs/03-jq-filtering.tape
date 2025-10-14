Output ../assets/demo-03-jq-filtering.gif

Require mozzy

Set FontSize 16
Set Width 1200
Set Height 700
Set Padding 20
Set Theme "Molokai"

Type "# Example 3: Inline JQ Filtering" Sleep 500ms Enter
Sleep 1s

Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --color"
Sleep 1s
Enter
Sleep 3s

Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .name --color"
Sleep 1s
Enter
Sleep 2s

Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .address.city --color"
Sleep 1s
Enter
Sleep 2s

Type "mozzy GET https://jsonplaceholder.typicode.com/users/1 --jq .company.name --color"
Sleep 1s
Enter
Sleep 3s

Type "# ✅ No external jq needed!" Sleep 500ms Enter
Sleep 2s
