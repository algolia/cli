package updateObjects

import (
	"encoding/json"
	"testing"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/stretchr/testify/assert"
)

func Test_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name       string
		json       string
		wantOut    ObjectsToUpdate
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "ObjectID, name/value",
			json: `
				[
					{
						"objectID": "23",
						"property1": "mj"
					}
				]`,
			wantOut: ObjectsToUpdate{
				{"objectID": "23", "property1": "mj"},
			},
			wantErr: false,
		},
		{
			name: "ObjectID, name/operation (increment)",
			json: `
				[
					{
						"objectID": "24",
						"property2": {
							"operation": "Increment",
							"value": 1
						}
					}
				]`,
			wantOut: ObjectsToUpdate{
				{
					"objectID": "24",
					"property2": search.PartialUpdateOperation{
						Operation: "Increment",
						Value:     1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "No ObjectID",
			json: `
				[
					{
						"property1": "mj"
					}
				]`,
			wantOut: ObjectsToUpdate{
				{
					"property1": "mj",
				},
			},
			wantErr: false,
		},
		{
			name: "Wrong operation",
			json: `
				[
					{
						"objectID": "24",
						"property2": {
							"operation": "Unknown",
							"value": "random"
						}
					}
				]`,
			wantErr:    true,
			wantErrMsg: "Invalid operation type for object 1",
		},
		{
			name: "Wrong operation structure",
			json: `
				[
					{
						"objectID": "24",
						"property2": {
							"type": "Unknown",
							"value": "random"
						}
					}
				]`,
			wantErr:    true,
			wantErrMsg: "Invalid operation for object 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var objectsToUpdate ObjectsToUpdate
			err := json.Unmarshal([]byte(tt.json), &objectsToUpdate)

			if err != nil {
				assert.Equal(t, tt.wantErr, true)
				assert.EqualError(t, err, tt.wantErrMsg)
				return
			}

			assert.EqualValues(t, objectsToUpdate, tt.wantOut)
		})
	}
}

func Test_isOperationTypeValid(t *testing.T) {
	tests := []struct {
		name          string
		operationType string
		wantOut       bool
	}{
		{
			name:          "Correct operation type",
			operationType: "Increment",
			wantOut:       true,
		},
		{
			name:          "Incorrect operation type (lowercase)",
			operationType: "increment",
			wantOut:       false,
		},
		{
			name:          "Incorrect operation type",
			operationType: "Unknown",
			wantOut:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantOut, isOperationTypeValid(tt.operationType))
		})
	}
}
