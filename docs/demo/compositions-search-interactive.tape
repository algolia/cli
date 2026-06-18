# Demo: `algolia compositions search --interactive`
#
# Showcases: confirm, multi-select parameter bag with type-to-filter, and a
# text input. Builds the search RequestBody, then runs the search and shows the
# results (the `demo` index holds recipe records).
#
# Record with:  vhs compositions-search-interactive.tape
# Requires:     vhs + ffmpeg, and `algolia` on PATH (sandbox profile).

Output compositions-search-interactive.gif

Set Shell "bash"
Set FontSize 22
Set Width 1650
Set Height 1000
Set Padding 10
Set TypingSpeed 60ms

Hide
Type "clear" Enter
Show

Type "algolia compositions search demo-comp --interactive"
Sleep 500ms
Enter
Sleep 1.5s

# 1) Set params? — confirm yes
Type "y"
Enter
Sleep 1s

# 2) params — multi-select bag; type to filter down to "query", toggle it, confirm
Type "query"
Sleep 1.5s
Space
Sleep 1s
Enter
Sleep 1.5s

# 3) params.query — text input
Type "cake"
Sleep 800ms
Enter

# build completes -> runs the search; show the recipe hits
Sleep 10s
