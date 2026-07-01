# Demo: parameter-bag deep dive in `compositions search --interactive`
#
# The search params are a large optional object, so the engine shows a
# multi-select. Use type-to-filter to pick a string + enum + int, then set each.
#
# Record with:  vhs compositions-search-parambag-interactive.tape
# Requires:     vhs + ffmpeg, `algolia` on PATH (sandbox profile).

Output compositions-search-parambag-interactive.gif

Set Shell "bash"
Set FontSize 18
Set Width 1650
Set Height 1100
Set Padding 10
Set TypingSpeed 60ms

Hide
Type "clear" Enter
Show

Type "algolia compositions search demo-comp --interactive"
Sleep 500ms
Enter
Sleep 1.5s

# 1) Set params? — yes
Type "y"
Enter
Sleep 1s

# 2) params multi-select with type-to-filter. Use a UNIQUE filter per field and
#    toggle the single top match (no arrow nav), clearing the filter between.
#    queryType (9 chars) -> toggle
Type "queryType"
Sleep 1.2s
Space
Sleep 600ms
Backspace@80ms 9
#    hitsPerPage (11 chars) -> toggle
Type "hitsPerPage"
Sleep 1.2s
Space
Sleep 800ms
Enter
Sleep 1.5s

# Selected fields are prompted in struct order: queryType, then hitsPerPage.
# 3) queryType — enum, validated
Type "prefixAll"
Enter
Sleep 1.2s

# 4) hitsPerPage — integer
Type "5"
Enter

# search runs -> show the hits
Sleep 10s
