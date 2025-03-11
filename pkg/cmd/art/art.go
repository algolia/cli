package art

import (
	"fmt"
	"time"

	"github.com/algolia/cli/pkg/auth"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

var art = `
           ________________________
         /                          \
        |            _____           |
        |       _.   XXXXX           |
        |      X/ .<&%%Y%%&>.        |
        |       .%#/   |##+\#%.      |
        |       %#` + "`" + `    |#/  ` + "`" + `#%      |
        |       %#     +     #%      |
        |       ` + "`" + `%\         /%` + "`" + `      |
        |        \#%b.___.d%#/       |
        |          ` + "`" + `` + "`" + `+===+` + "`" + `` + "`" + `         |
        |                            |
         \ ________________________ /
		 
      Congratulations on your epic search!
             You found the ☼art☼!
`

func NewArtCmd(f *cmdutil.Factory) *cobra.Command {
	loadingMessages := []string{
		"The legends speak of an Algolia developer so strong...",
		"so capable...",
		"so powerful at search and discovery...",
		"that they could find the hidden art at the root of all things...",
		"So spake the legends...",
	}

	cmd := &cobra.Command{
		Use:   "art",
		Short: "We've been searching for the art for so long...!",
		Run: func(cmd *cobra.Command, args []string) {
			io := f.IOStreams
			io.StartProgressIndicatorWithLabel("LEGENDARY QUEST ACCEPTED")
			time.Sleep(2 * time.Second)
			for i := 0; i < len(loadingMessages); i++ {
				io.UpdateProgressIndicatorLabel(loadingMessages[i])
				time.Sleep(2 * time.Second)
			}
			io.StopProgressIndicator()
			fmt.Println(art)
		},
		Hidden: true,
	}

	auth.DisableAuthCheck(cmd)

	return cmd
}
