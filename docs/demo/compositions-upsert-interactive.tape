# Demo: `algolia compositions upsert --interactive`
#
# Showcases every prompt the engine produces:
#   text input, optional-skip, single-select (oneOf union), confirm,
#   multi-select (parameter bag), bool, integer, enum, slice, map.
#
# Record with:  vhs compositions-upsert-interactive.tape
# Requires:     vhs + ffmpeg, and `algolia` on PATH.
# NOTE: the command performs a real PutComposition at the end — run it against a
#       throwaway / sandbox application profile.
#
# Prompt order is tied to the SDK Composition struct; if the SDK changes, replay
# and adjust the keystrokes (see pkg/cmd/compositions trace approach in git history).

Output compositions-upsert-interactive.gif

Set Shell "bash"
Set FontSize 16
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

# 1) name — required text input
Type "Demo Composition"
Sleep 400ms
Enter
Sleep 800ms

# 2) description — optional, press Enter to skip
Sleep 500ms
Enter
Sleep 800ms

# 3) behavior — single-select union; keep the highlighted first option (injection)
Sleep 600ms
Enter
Sleep 800ms

# 4) Set ...injection.main.source? — confirm yes
Type "y"
Enter
Sleep 1s

# 5) source — single-select union; arrow down to the search source
Down
Sleep 2s
Enter
Sleep 1s

# 6) search.index — required text input
Type "demo"
Enter
Sleep 800ms

# 7) Set ...search.params? — confirm yes
Type "y"
Enter
Sleep 1s

# 8) params — multi-select parameter bag; toggle the first two
#    (advancedSyntax, advancedSyntaxFeatures), then confirm
Space
Sleep 400ms
Down
Space
Sleep 600ms
Enter
Sleep 2s

# 9) advancedSyntax — bool entered as true/false
Type "true"
Enter
Sleep 1000ms

# 10) how many advancedSyntaxFeatures — integer (count)
Type "1"
Enter
Sleep 2200ms

# 11) advancedSyntaxFeatures[0] — enum, validated against the SDK
Type "exactPhrase"
Enter
Sleep 2200ms

# 12) how many injectedItems — integer; 0 skips the slice
Type "0"
Enter
Sleep 2200ms

# 13) Set ...deduplication? — confirm no (press Enter for the default)
Enter
Sleep 800ms

# 14) Set sortingStrategy? — confirm yes
Type "y"
Enter
Sleep 800ms

# 15) how many sortingStrategy entries — integer (count)
Type "1"
Enter
Sleep 600ms

# 16) sortingStrategy key[0] — map key
Type "Price (asc)"
Enter
Sleep 600ms

# 17) sortingStrategy value — map value
Type "demo2"
Sleep 2s
Enter

Enter
Sleep 10s
