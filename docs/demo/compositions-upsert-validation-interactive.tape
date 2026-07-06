# Demo: validation + retry in `compositions upsert --interactive`
#
# The engine never fails-fast: a bad entry re-prompts in place and prior answers
# are kept. Shows required-empty, bad-integer, and invalid-enum retries.
#
# Record with:  vhs compositions-upsert-validation-interactive.tape
# Requires:     vhs + ffmpeg, `algolia` on PATH (sandbox profile; real PutComposition).

Output compositions-upsert-validation-interactive.gif

Set Shell "bash"
Set FontSize 18
Set Width 1650
Set Height 1000
Set Padding 10
Set TypingSpeed 60ms

Hide
Type "clear" Enter
Show

Type "algolia compositions upsert demo-comp --interactive"
Sleep 500ms
Enter
Sleep 1.5s

# 1) name — required; submit empty first to show the re-prompt, then fill it
Enter
Sleep 1.5s
Type "Demo Composition"
Sleep 1s
Enter
Sleep 800ms

# 2) description — optional, skip
Enter
Sleep 600ms

# 3) behavior — single-select union; keep injection
Enter
Sleep 800ms

# 4) Set ...injection.main.source? — yes
Type "y"
Sleep 1s
Enter
Sleep 800ms

# 5) source — single-select union; pick the search source
Down
Sleep 600ms
Enter
Sleep 800ms

# 6) search.index — required; submit empty first, then fill it
Enter
Sleep 1.5s
Type "demo"
Sleep 1s
Enter
Sleep 800ms

# 7) Set ...search.params? — no
Enter
Sleep 800ms

# 8) how many injectedItems — bad value first ("abc"), then 0
Type "abc"
Sleep 1s
Enter
Sleep 1.8s
Type "0"
Sleep 1s
Enter
Sleep 800ms

# 9) Set ...deduplication? — yes
Type "y"
Sleep 1s
Enter
Sleep 800ms

# 10) deduplication.positioning — invalid enum first, then a valid value
Type "bogus"
Sleep 1s
Enter
Sleep 2s
Type "highest"
Sleep 1s
Enter
Sleep 800ms

# 11) Set sortingStrategy? — no
Enter

# build completes -> PutComposition + wait for published
Sleep 10s
