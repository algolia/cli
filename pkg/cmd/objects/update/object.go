package update

import (
	"encoding/json"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/utils"
	"github.com/mitchellh/mapstructure"
)

// Object is a map[string]interface{} that can be unmarshalled from a JSON object
// The object must have an objectID field
// Each field could be either an `search.PartialUpdateOperation` or a scalar value
type Object map[string]interface{}

// Valid operations
const (
	Increment     string = "Increment"
	Decrement     string = "Decrement"
	Add           string = "Add"
	AddUnique     string = "AddUnique"
	IncrementSet  string = "IncrementSet"
	IncrementFrom string = "IncrementFrom"
)

// ValidateOperation checks that the operation is valid
func ValidateOperation(p search.PartialUpdateOperation) error {
	allowedOperations := []string{Increment, Decrement, Add, AddUnique, IncrementSet, IncrementFrom}
	extra := fmt.Sprintf("valid operations are %s", utils.SliceToReadableString(allowedOperations))

	if p.Operation == "" {
		return fmt.Errorf("missing operation")
	}
	if !utils.Contains(allowedOperations, p.Operation) {
		return fmt.Errorf("invalid operation \"%s\" (%s)", p.Operation, extra)
	}
	return nil
}

// UnmarshalJSON unmarshals a JSON object into an Object
func (o *Object) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	// The object must be a map[string]interface{}
	switch v := v.(type) {
	case map[string]interface{}:
		*o = v
	default:
		return fmt.Errorf("invalid object: %v", v)
	}

	// The object must have an objectID
	if _, ok := (*o)["objectID"]; !ok {
		return fmt.Errorf("objectID is required")
	}

	// Each field could be either an `search.PartialUpdateOperation` or a scalar value
	for k, v := range *o {
		switch v := v.(type) {
		case map[string]interface{}:
			var op search.PartialUpdateOperation
			if err := mapstructure.Decode(v, &op); err != nil {
				return err
			}
// Check that the operation is valid
			if err := ValidateOperation(op); err != nil {
				return err
			}
			(*o)[k] = op
		}
	}

	return nil
}
