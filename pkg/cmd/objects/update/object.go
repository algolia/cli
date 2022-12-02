package updateObjects

import (
	"encoding/json"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/algolia/cli/pkg/utils"
)

type ObjectsToUpdate []map[string]interface{}

func (o *ObjectsToUpdate) UnmarshalJSON(data []byte) error {
	var rawObjectsToUpdate []map[string]interface{}
	err := json.Unmarshal(data, &rawObjectsToUpdate)
	if err != nil {
		return err
	}

	var objectsToUpdate ObjectsToUpdate
	for objectIndex := range rawObjectsToUpdate {
		rawObject := rawObjectsToUpdate[objectIndex]
		objectToUpdate := map[string]interface{}{}

		for name, rawValue := range rawObject {
			valueString, isValueString := rawValue.(string)
			if name != "objectID" {
				_, isValueBool := rawValue.(bool)
				_, isValueFloat := rawValue.(float64)

				if isValueBool || isValueFloat || isValueString {
					objectToUpdate[name] = rawValue
				} else {
					mapValue, isMap := rawValue.(map[string]interface{})
					if !isMap {
						return fmt.Errorf("Object value/operation not recognized for object %d", objectIndex+1)
					}
					mapValueString, isMapValueString := mapValue["value"].(string)
					// JSON unmarshal numbers to float64 by default
					mapValueFloat, isMapValueFloat := mapValue["value"].(float64)
					mapValueInt := int(mapValueFloat)

					operationString, ok := mapValue["operation"].(string)
					if !ok {
						return fmt.Errorf("Invalid operation for object %d", objectIndex+1)
					}
					isOperationValid := isOperationTypeValid(operationString)
					if !isOperationValid {
						return fmt.Errorf("Invalid operation type for object %d", objectIndex+1)
					}

					var operationValue interface{}
					if isMapValueFloat {
						operationValue = mapValueInt
					} else if isMapValueString {
						operationValue = mapValueString
					}
					objectToUpdate[name] = search.PartialUpdateOperation{
						Operation: operationString,
						Value:     operationValue,
					}
				}
			} else if !isValueString {
				return fmt.Errorf("Error at object %d: objectID should be a string", objectIndex+1)
			} else {
				objectToUpdate["objectID"] = valueString
			}
		}
		objectsToUpdate = append(objectsToUpdate, objectToUpdate)
	}
	*o = objectsToUpdate
	return nil
}

func isOperationTypeValid(value string) bool {
	return utils.Contains([]string{"Increment", "Decrement", "Add", "Remove",
		"AddUnique", "IncrementFrom", "IncrementSet"}, value)
}
