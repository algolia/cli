package art

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/algolia/cli/pkg/cmdutil"

)

var artString = `

           ________________________
         /                          \
        |            _____           |
        |       _.   XXXXX           |
        |      X/ .<&%%Y%%&>.        |
        |       .%#/   |##+\#%.      |
        |      .%#°    |#/  °#%      |
        |       %#     °     #%      |
        |       °%\         /%°      |
        |        \#%b.___.d%#/       |
        |           °+===+°          |
        |                            |
         \ ________________________ /
		 
      Congratulations on your epic search!
             You found the art!
`



func NewArtCmd(f *cmdutil.Factory) *cobra.Command {
	// artCmd represents the art command
	var artCmd = &cobra.Command{
		Use:   "art",
		Short: "Wow, we found the art!",
		Long: `The legends speak of a search API so strong...
		so powerful...
		that it could find the root of all things...
		So the legends speak...`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(artString)
		},
	}

	return artCmd
}