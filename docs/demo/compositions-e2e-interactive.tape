# Demo: end-to-end story
#   1. build a composition interactively
#   2. get it back to show it landed
#   3. search it interactively (the `demo` index holds recipes)
#
# Record with:  vhs compositions-e2e-interactive.tape
# Requires:     vhs + ffmpeg, `algolia` on PATH (sandbox profile; real API calls).

Output compositions-e2e-interactive.gif

Set Shell "bash"
Set FontSize 22
Set Width 1650
Set Height 1100
Set Padding 10
Set TypingSpeed 60ms

Hide
Type "clear" Enter
Show

# ── 1. upsert interactively (minimal valid composition) ───────────────────────
Type "algolia compositions upsert demo-comp --interactive"
Sleep 500ms
Enter
Sleep 1.5s

Type "Demo Composition"
Sleep 600ms
Enter
Sleep 600ms
Enter
Sleep 600ms
# behavior union -> injection
Enter
Sleep 600ms
# Set source? yes
Type "y"
Enter
Sleep 600ms
# source union -> search
Down
Sleep 1s
Enter
Sleep 600ms
# index
Type "demo"
Enter
Sleep 800ms
# Set params? no
Enter
Sleep 600ms
# how many injectedItems -> 0
Type "0"
Sleep 800ms
Enter
Sleep 600ms
# Set deduplication? no
Enter
Sleep 1s
# Set sortingStrategy? no
Enter
Sleep 6s

# ── 2. get it back ────────────────────────────────────────────────────────────
Type "algolia compositions get demo-comp"
Sleep 400ms
Enter
Sleep 5s

# ── 3. search it interactively ────────────────────────────────────────────────
Type "algolia compositions search demo-comp --interactive"
Sleep 400ms
Enter
Sleep 1.5s
# Set params? yes
Type "y"
Sleep 1s
Enter
Sleep 1s
# multi-select: filter to "query", toggle, confirm
Type "query"
Sleep 1.2s
Space
Sleep 800ms
Enter
Sleep 1.2s
# query value
Type "cake"
Sleep 4s
Enter

# search runs -> show the recipe hits
Sleep 10s
