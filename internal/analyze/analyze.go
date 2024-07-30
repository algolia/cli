package analyze

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v4/algolia/search"
)

// AttributeType is an enum for the different types of attributes.
type AttributeType string

const (
	// String is the type for string attributes.
	String AttributeType = "string"
	// Numeric is the type for numeric attributes (integers and floats).
	Numeric AttributeType = "numeric"
	// Boolean is the type for boolean attributes.
	Boolean AttributeType = "boolean"
	// Array is the type for array attributes.
	Array AttributeType = "array"
	// Object is the type for object attributes.
	Object AttributeType = "object"
	// Null is the type for null attributes.
	Null AttributeType = "null"
	// Undefined is the type for undefined attributes (not present in the object).
	Undefined AttributeType = "undefined"
)

// AttributeStats contains the stats for a single attribute of an Algolia object.
type AttributeStats struct {
	Count      int                       `json:"count"`
	Percentage float64                   `json:"percentage"`
	Types      map[AttributeType]float64 `json:"types"`
	InSettings []string                  `json:"inSettings,omitempty"`

	Values map[interface{}]int `json:"values,omitempty"`
}

// Stats contains the stats for an Algolia index.
type Stats struct {
	TotalRecords int
	Attributes   map[string]*AttributeStats
}

// ComputeStats computes the stats for the given index.
func ComputeStats(
	records []search.Hit,
	settings search.SettingsResponse,
	limit int,
	only string,
	counter chan int,
) (*Stats, error) {
	settingsMap := settingsAsMap(settings)
	stats := &Stats{
		Attributes: make(map[string]*AttributeStats),
	}

	for _, r := range records {
		if limit > 0 && stats.TotalRecords >= limit {
			break
		}
		stats.TotalRecords++
		counter <- 1
		stats = computeObjectStats(stats, "", r.AdditionalProperties, only)
	}

	for key, value := range stats.Attributes {
		// Calculate the percentage of each attribute rounded up
		if !strings.Contains(key, ".") {
			value.Percentage = float64(value.Count) * 100 / float64(stats.TotalRecords)
		} else {
			// If the attribute is a nested one, compute the percentage based on the parent attribute count
			value.Percentage = float64(value.Count) * 100 / float64(stats.Attributes[parentKey(key)].Count)
		}

		// Calculate the percentage of each type rounded up
		for typeKey, typeValue := range value.Types {
			if !strings.Contains(key, ".") {
				value.Types[typeKey] = float64(int(typeValue*100)) / float64(stats.TotalRecords)
			} else {
				// If the attribute is a nested one, compute the percentage based on the parent attribute count
				value.Types[typeKey] = float64(int(typeValue*100)) / float64(stats.Attributes[parentKey(key)].Count)
			}
		}

		// Leftover attributes are "undefined"
		value.Types[Undefined] = 100 - value.Types[String] - value.Types[Numeric] - value.Types[Boolean] - value.Types[Array] - value.Types[Object] - value.Types[Null]
		if value.Types[Undefined] <= 0 {
			delete(value.Types, Undefined)
		}

		// Add the settings where the attribute is present
		value.InSettings = inSettings(settingsMap, key)

		stats.Attributes[key] = value
	}

	return stats, nil
}

// computeObjectStats computes the stats for the given object.
func computeObjectStats(s *Stats, p string, o map[string]interface{}, only string) *Stats {
	for key, value := range o {
		if only != "" && only != key {
			continue
		}

		var fullPath string
		if p == "" {
			fullPath = key
		} else {
			fullPath = fmt.Sprintf("%s.%s", p, key)
		}

		if getType(value) == Object {
			v, ok := value.(map[string]interface{})
			if ok {
				s = computeObjectStats(s, fullPath, v, only)
			}
		}

		if _, ok := s.Attributes[fullPath]; !ok {
			s.Attributes[fullPath] = &AttributeStats{
				Types:  make(map[AttributeType]float64),
				Values: make(map[interface{}]int),
			}
		}
		s.Attributes[fullPath].Count++
		s.Attributes[fullPath].Types[getType(value)]++

		if only != "" {
			switch getType(value) {
			case String, Numeric, Boolean:
				s.Attributes[fullPath].Values[value]++
			case Array:
				for _, v := range value.([]interface{}) {
					s.Attributes[fullPath].Values[v]++
				}
			}
		}
	}

	return s
}

// getType returns the type of the given value
func getType(value interface{}) AttributeType {
	switch value.(type) {
	case string:
		return String
	case int, float64:
		return Numeric
	case bool:
		return Boolean
	case []interface{}:
		return Array
	case map[string]interface{}:
		return Object
	default:
		return Null
	}
}

// settingsAsMap converts the given settings to a map.
// We marshal and unmarshal the settings to avoid having to write the conversion code ourselves.
func settingsAsMap(s search.SettingsResponse) map[string]interface{} {
	var settingsMap map[string]interface{}
	var settingsBytes []byte
	settingsBytes, err := s.MarshalJSON()
	if err != nil {
		return nil
	}
	err = json.Unmarshal(settingsBytes, &settingsMap)
	if err != nil {
		return nil
	}
	return settingsMap
}

// inSettings returns a slice of strings containing the settings where the given key is present
// TODO: Handle multiple attributes like `firstAttribute,secondAttribute` at the same level
func inSettings(s map[string]interface{}, key string) []string {
	var result []string
	possiblePatterns := []string{
		"%s",
		"ordered(%s)",
		"unordered(%s)",
		"searchable(%s)",
		"exact(%s)",
		"filterOnly(%s)",
		"afterDistinct(%s)",
		"desc(%s)",
	}
	for s, v := range s {
		if v, ok := v.([]interface{}); ok {
			for _, v := range v {
				var toSearch []string
				for _, pattern := range possiblePatterns {
					toSearch = append(toSearch, fmt.Sprintf(pattern, key))
				}

				for _, pattern := range toSearch {
					if fmt.Sprintf("%s", v) == pattern {
						result = append(result, s)
						break
					}
				}
			}
		}
	}

	return result
}

// parentKey returns the parent key of the given key
func parentKey(key string) string {
	return strings.Join(strings.Split(key, ".")[:len(strings.Split(key, "."))-1], ".")
}
