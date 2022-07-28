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
             You found the ☼art☼!
`



func NewArtCmd(f *cmdutil.Factory) *cobra.Command {
	// artCmd represents the art command
	var artCmd = &cobra.Command{
		Use:   "art",
		Short: "We've been searching for the art for so long...!",
		Long: `LEGENDARY QUEST ACCEPTED:
		The legends speak of an Algolia developer so strong...
		so capable...
		so powerful at search and discovery...
		that they could find the hidden art at the root of all things...
		So spake the legends...`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(artString)
		},
		Hidden:  true,
	}

	return artCmd
}