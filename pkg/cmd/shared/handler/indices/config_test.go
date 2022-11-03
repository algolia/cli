package config

import (
	"testing"

	"github.com/algolia/cli/pkg/iostreams"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateExportConfigFlags(t *testing.T) {
	tests := []struct {
		name        string
		opts        ExportOptions
		wantsErr    bool
		wantsErrMsg string
	}{
		{
			name: "No existing indice",
			opts: ExportOptions{
				Indices:         []string{"INDICE_1"},
				Scope:           []string{"settings", "rules", "synonyms"},
				ExistingIndices: []string{},
			},
			wantsErr:    true,
			wantsErrMsg: "X Indice 'INDICE_1' doesn't exist",
		},
		{
			name: "No scope",
			opts: ExportOptions{
				Indices:         []string{"INDICE_1"},
				Scope:           []string{},
				ExistingIndices: []string{"INDICE_1", "INDICE_2"},
			},
			wantsErr:    true,
			wantsErrMsg: "X required flag scope not set",
		},
		{
			name: "Full scope with existing indices",
			opts: ExportOptions{
				Indices:         []string{"INDICE_1"},
				Scope:           []string{"settings", "rules", "synonyms"},
				ExistingIndices: []string{"INDICE_1", "INDICE_2"},
			},
			wantsErr: false,
		},
		{
			name: "Full score, existing indices with directory",
			opts: ExportOptions{
				Indices:         []string{"INDICE_1"},
				Scope:           []string{"settings", "rules", "synonyms"},
				ExistingIndices: []string{"INDICE_1", "INDICE_2"},
				Directory:       "test/folder",
			},
			wantsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			tt.opts.IO = io

			err := ValidateExportConfigFlags(tt.opts)

			if tt.wantsErr {
				assert.EqualError(t, err, tt.wantsErrMsg)
				return
			}
		})
	}
}
