package tools

import (
	"encoding/json"
	"fmt"
	"strings"
)

const toolTypeAlgoliaSearchIndex = "algolia_search_index"

func mergeAlgoliaSearchIndexTool(
	toolsJSON []byte,
	toolName string,
	indexName string,
	description string,
	searchParameters []byte,
) ([]byte, error) {
	toolName = strings.TrimSpace(toolName)
	if err := validateToolName(toolName); err != nil {
		return nil, err
	}
	if strings.TrimSpace(indexName) == "" {
		return nil, fmt.Errorf("index name must not be empty")
	}

	arr, err := normalizeToolsArray(toolsJSON)
	if err != nil {
		return nil, err
	}

	entry := map[string]any{
		"index":       indexName,
		"description": description,
	}
	if len(searchParameters) > 0 {
		trim := strings.TrimSpace(string(searchParameters))
		if trim != "" {
			var sp any
			if err := json.Unmarshal(searchParameters, &sp); err != nil {
				return nil, fmt.Errorf("invalid --search-parameters JSON: %w", err)
			}
			entry["searchParameters"] = sp
		}
	}

	for i, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if typ, _ := m["type"].(string); typ != toolTypeAlgoliaSearchIndex {
			continue
		}
		if _, has := m["name"]; !has || strings.TrimSpace(fmt.Sprint(m["name"])) == "" {
			m["name"] = toolName
		}
		indices, _ := m["indices"].([]any)
		for _, ex := range indices {
			em, ok := ex.(map[string]any)
			if !ok {
				continue
			}
			if exIdx, _ := em["index"].(string); exIdx == indexName {
				return nil, fmt.Errorf(
					"index %q is already present on this agent's algolia_search_index tool",
					indexName,
				)
			}
		}
		m["indices"] = append(indices, entry)
		arr[i] = m
		return json.Marshal(arr)
	}

	newTool := map[string]any{
		"name":    toolName,
		"type":    toolTypeAlgoliaSearchIndex,
		"indices": []any{entry},
	}
	arr = append(arr, newTool)
	return json.Marshal(arr)
}

func normalizeToolsArray(toolsJSON []byte) ([]any, error) {
	if len(toolsJSON) == 0 {
		return []any{}, nil
	}
	trim := strings.TrimSpace(string(toolsJSON))
	if trim == "" || trim == "null" {
		return []any{}, nil
	}

	var arr []any
	if err := json.Unmarshal([]byte(trim), &arr); err != nil {
		return nil, fmt.Errorf("agent tools must be a JSON array (or empty): %w", err)
	}
	return arr, nil
}

// validateToolName enforces AlgoliaSearchToolConfig-Input.name (OpenAPI:
// minLength 3, maxLength 32).
func validateToolName(name string) error {
	n := len([]rune(name))
	switch {
	case n < 3:
		return fmt.Errorf("tool name must be at least 3 characters (OpenAPI Agent Studio schema)")
	case n > 32:
		return fmt.Errorf("tool name must be at most 32 characters (OpenAPI Agent Studio schema)")
	}
	return nil
}
