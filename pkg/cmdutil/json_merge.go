package cmdutil

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/iostreams"
)

// MergeFileAndFlagsInto reads an optional JSON `file` (path or "-" for stdin),
// then overlays the named cobra flag values on top, then unmarshals the
// merged map into out. Flag values win over file values for the same key.
//
// Used by commands that accept both `-F file.json` and a generated flag set
// derived from an OpenAPI schema (see e.g. AgentConfigCreate flags).
func MergeFileAndFlagsInto(io *iostreams.IOStreams, file string, cmd *cobra.Command, flagNames []string, out any) error {
	merged := map[string]any{}
	if file != "" {
		b, err := ReadFile(file, io.In)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, &merged); err != nil {
			return fmt.Errorf("parsing %s: %w", file, err)
		}
	}

	flagVals, err := FlagValuesMap(cmd.Flags(), flagNames...)
	if err != nil {
		return err
	}
	for k, v := range flagVals {
		merged[k] = v
	}

	tmp, err := json.Marshal(merged)
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp, out)
}
