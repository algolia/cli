# Demo: `algolia compositions rules upsert --interactive`
#
# Showcases: slice (conditions, tags), enum (anchoring), optional-skip,
# single-select unions (consequence.behavior + injection source), confirm,
# optional bool (enabled). A main source is selected because the composition
# API rejects an injection behavior without one (SDK quirk, not a bug).
#
# Record with:  vhs compositions-rules-upsert-interactive.tape
# Requires:     vhs + ffmpeg, and `algolia` on PATH.
# NOTE: performs a real PutCompositionRule at the end — use a sandbox profile.

Output compositions-rules-upsert-interactive.gif

Set Shell "bash"
Set FontSize 16
Set Width 1650
Set Height 1000
Set Padding 10
Set TypingSpeed 60ms

Hide
Type "clear" Enter
Show

Type "algolia compositions rules upsert demo-comp rule-1 --interactive"
Sleep 500ms
Enter
Sleep 1.5s

# 1) how many conditions — integer (count)
Type "1"
Enter
Sleep 1s

# 2) conditions[0].pattern — optional text
Type "cake"
Enter
Sleep 800ms

# 3) conditions[0].anchoring — enum, validated against the SDK
Type "contains"
Enter
Sleep 800ms

# 4) conditions[0].context — optional, Enter to skip
Enter
Sleep 500ms
# 5) conditions[0].filters — optional, Enter to skip
Enter
Sleep 500ms
# 6) conditions[0].sortBy — optional, Enter to skip
Enter
Sleep 800ms

# 7) consequence.behavior — single-select union; keep injection (highlighted)
Sleep 600ms
Enter
Sleep 800ms

# 8) Set ...injection.main.source? — confirm yes (a source is required by the API)
Type "y"
Enter
Sleep 800ms

# 8a) source — single-select union; arrow down to the search source
Down
Sleep 800ms
Enter
Sleep 800ms

# 8b) ...search.index — required text input
Type "demo"
Enter
Sleep 800ms

# 8c) Set ...search.params? — confirm no (Enter for default)
Enter
Sleep 800ms

# 9) how many injectedItems — integer; 0 to skip
Type "0"
Enter
Sleep 800ms

# 10) Set ...deduplication? — confirm no (Enter for default)
Enter
Sleep 800ms

# 11) description — optional, Enter to skip
Enter
Sleep 600ms

# 12) enabled — optional bool entered as true/false
Type "true"
Enter
Sleep 800ms

# 13) how many validity — integer; 0 to skip
Type "0"
Enter
Sleep 800ms

# 14) how many tags — integer (count)
Type "2"
Enter
Sleep 600ms

# 15) tags[0]
Type "sale"
Enter
Sleep 500ms

# 16) tags[1]
Type "seasonal"
Sleep 1s
Enter

# build completes -> PutCompositionRule + wait for published
Sleep 10s
