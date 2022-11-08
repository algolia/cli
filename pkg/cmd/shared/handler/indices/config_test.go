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

func Test_ValidateImportConfigFlags(t *testing.T) {
	tests := []struct {
		name        string
		opts        ImportOptions
		wantsErr    bool
		wantsErrMsg string
	}{
		{
			name: "Import rules",
			opts: ImportOptions{
				Scope:    []string{"rules"},
				FilePath: "config_mock.json",
			},
			wantsErr: false,
		},
		{
			name: "Import rules and synonyms",
			opts: ImportOptions{
				Scope:    []string{"rules", "synonyms"},
				FilePath: "config_mock.json",
			},
			wantsErr: false,
		},
		{
			name: "Clear existing rules without rules in scope",
			opts: ImportOptions{
				Scope:              []string{"synonyms"},
				FilePath:           "config_mock.json",
				ClearExistingRules: true,
			},
			wantsErr:    true,
			wantsErrMsg: "X Cannot clear existing rules if rules are not in scope",
		},
		{
			name: "Clear existing synonyms without synonyms in scope",
			opts: ImportOptions{
				Scope:                 []string{"rules"},
				FilePath:              "config_mock.json",
				ClearExistingSynonyms: true,
			},
			wantsErr:    true,
			wantsErrMsg: "X Cannot clear existing synonyms if synonyms are not in scope",
		},
		{
			name: "Wrong file path",
			opts: ImportOptions{
				Scope:    []string{"settings", "rules", "synonyms"},
				FilePath: "wrong_path.json",
			},
			wantsErr:    true,
			wantsErrMsg: "X An error occured when opening file: open wrong_path.json: no such file or directory",
		},
		{
			name: "Import settings",
			opts: ImportOptions{
				Scope:    []string{"settings"},
				FilePath: "config_mock.json",
			},
			wantsErr:    true,
			wantsErrMsg: "X No setting found in config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			tt.opts.IO = io

			err := ValidateImportConfigFlags(&tt.opts)

			if tt.wantsErr {
				assert.EqualError(t, err, tt.wantsErrMsg)
				return
			}
			assert.Equal(t, nil, err)
		})
	}
}
