# Algolia CLI Tool: Compositions Interactive Mode

## What is it?

The Algolia CLI (`algolia`) is Algolia's official command-line tool for managing your application from the terminal: indices, records, rules, API keys, synonyms, and Compositions.

- Official documentation: https://www.algolia.com/doc/tools/cli/get-started/overview/
- Source code (public): https://github.com/algolia/cli

## How?

This page covers **Compositions Interactive Mode**. It lets you build a composition request body by answering prompts, instead of hand-writing JSON. Add `--interactive` (or `-i`) to a compositions command and the CLI walks you through the body field by field.

## Why use it

- Composition bodies are large and deeply nested (behavior variants, injection sources, dozens of search parameters). Interactive mode means you do not need to know the JSON shape.
- You only answer for the fields you care about. Everything optional can be skipped.
- Input is checked as you type. A bad value is rejected on the spot and re-asked, so you never lose the answers you already gave.

## Requirements

- An interactive terminal. `-i` needs a real TTY. In scripts or CI, keep using `--file`.
- A profile with the right permissions: `editSettings` for `upsert` and `rules upsert`, `search` for `search`.

## Supported commands

| Command | Builds | Demo |
|---------|--------|------|
| `algolia compositions upsert <id> -i` | a composition | see below |
| `algolia compositions rules upsert <id> <rule-id> -i` | a composition rule | see below |
| `algolia compositions search <id> -i` | a search request (the `<query>` argument becomes optional) | see below |

### Upsert a composition

```
# Build a composition interactively
$ algolia compositions upsert my-comp --interactive

# Still supported: from a JSON file or stdin
$ algolia compositions upsert my-comp --file body.json
$ cat body.json | algolia compositions upsert my-comp --file -
```

![Upsert a composition interactively](../compositions-upsert-interactive.gif)

### Upsert a composition rule

```
# Build a rule interactively
$ algolia compositions rules upsert my-comp rule-1 --interactive

# Still supported: from a JSON file or stdin
$ algolia compositions rules upsert my-comp rule-1 --file rule.json
$ cat rule.json | algolia compositions rules upsert my-comp rule-1 --file -
```

![Upsert a rule interactively](../compositions-rules-upsert-interactive.gif)

### Search a composition

```
# Build the search request interactively
$ algolia compositions search my-comp --interactive

# Still supported: query and flags
$ algolia compositions search my-comp "running shoes"
$ algolia compositions search my-comp "shirt" --filters "brand:Nike"
$ algolia compositions search my-comp "shirt" --hits-per-page 20 --page 2
```

![Search a composition interactively](../compositions-search-interactive.gif)

## How the prompts work

You will see a few kinds of prompt:

- **Text.** Type a value and press Enter.
- **Yes/no.** For example `Set ...source?`. Type `y` or press Enter for no.
- **Pick one.** For "one of" choices such as the composition `behavior` or the injection source. Use the arrow keys and Enter.
- **Pick many.** For large parameter objects (the search params). Type to filter the list, press Space to toggle an item, Enter to confirm.

![Filtering and selecting search parameters](../compositions-search-parambag-interactive.gif)

- **Optional fields** can be skipped by pressing Enter on an empty line.
- **Lists and maps** first ask how many items you want, then prompt for each one.
- **Enums** only accept allowed values. An invalid value is rejected and re-asked.
- **Validation and retry.** An empty required field, a non-number where a number is expected, or an invalid enum is rejected in place and re-asked. Your earlier answers are kept.

![Invalid entries are re-prompted](../compositions-upsert-validation-interactive.gif)

- Press `Ctrl-C` at any time to abort without sending anything.

## End to end

Build a composition, read it back, then search it:

```
$ algolia compositions upsert demo-comp --interactive
$ algolia compositions get demo-comp
$ algolia compositions search demo-comp --interactive
```

![Build, get, then search](../compositions-e2e-interactive.gif)

## Tips and gotchas

- An injection behavior needs a source. Pick one when prompted (for example the search source and its index). The API rejects an injection with no source.
    - If in doubt, use the cli tool to check your available indices:
    ```sh
    algolia indices list
    ```
- `-i` needs a terminal. Use `--file` for automation and CI.
- The prompts follow the Algolia SDK model for compositions, so the flow stays in sync as the API changes.
- When in doubt, consult the official API reference documentation for description of each option and arguemt: 
    - https://www.algolia.com/doc/rest-api/composition.
- Per-command help: `algolia compositions upsert --help`, `algolia compositions rules upsert --help`, `algolia compositions search --help`
