# Demo: `algolia compositions sorting-strategy --interactive`
#
# Builds the sortingStrategy map (sort label -> index/replica) via prompts:
# a count, then a label + index for each entry. The sort index must be related
# to the composition's main index; here the main index `demo` is used so the
# update is accepted (in practice you point labels at replicas of the main index).
#
# Record with:  vhs compositions-sorting-strategy-interactive.tape
# Requires:     vhs + ffmpeg, `algolia` on PATH (sandbox profile; real API call).

Output compositions-sorting-strategy-interactive.gif

Set Shell "bash"
Set FontSize 24
Set Width 1650
Set Height 900
Set Padding 10
Set TypingSpeed 60ms

Hide
Type "clear" Enter
Show

Type "algolia compositions sorting-strategy demo-comp --interactive"
Sleep 500ms
Enter
Sleep 1.5s

# 1) how many entries — integer (count)
Type "2"
Sleep 1s
Enter
Sleep 1s

# 2) entry 0 — sort label, then the index that implements it
Type "Price (asc)"
Sleep 1s
Enter
Sleep 800ms
Type "demo"
Sleep 1s
Enter
Sleep 1s

# 3) entry 1 — a replica of the main index
Type "Price (desc)"
Sleep 1s
Enter
Sleep 800ms
Type "demo_replica"
Sleep 1s
Enter

# update runs -> POST sortingStrategy, wait for the task, show the taskID
Sleep 7s

Type "algolia compositions get demo-comp"
Sleep 500ms
Enter
Sleep 10s
