package shared

import "testing"

func Test_GetSynonymSuccessMessage(t *testing.T) {
	tests := []struct {
		name         string
		synonymFlags SynonymFlags
		saveOptions  SaveOptions
		wantsOutput  string
		saveWording  string
	}{
		{
			name: "Create regular synonym",
			synonymFlags: SynonymFlags{
				SynonymID: "23",
				Synonyms:  []string{"mj", "goat"},
			},
			saveOptions: SaveOptions{
				Indice: "legends",
			},
			wantsOutput: "",
			saveWording: "created",
		},
		// {
		// 	name: "Create one way synonym",
		// 	saveOptions: SaveOptions{
		// 		SynonymID: "23",
		// 		Synonyms:  []string{"mj", "goat"}},
		// 	wantsOutput: "",
		// 	saveWording: "created",
		// },
		// {
		// 	name: "Create regular synonym",
		// 	saveOptions: SaveOptions{
		// 		SynonymID: "23",
		// 		Synonyms:  []string{"mj", "goat"}},
		// 	wantsOutput: "",
		// 	saveWording: "created",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// outputMessage := GetSynonymSuccessWording(tt.saveOptions, tt.saveWording)

			// assert.Equal(t, outputMessage, tt.wantsOutput)
		})
	}
}
