package status

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/iostreams"
)

type StatusOptions struct {
	Config config.IConfig
	IO     *iostreams.IOStreams
	Status StatusResponse
}

type StatusResponse struct {
	Status Status
}

// TODO: use generation for this
type Status struct {
	C1Br    string `json:"c1-br"`
	C1Ca    string `json:"c1-ca"`
	C1De    string `json:"c1-de"`
	C1Hk    string `json:"c1-hk"`
	C1In    string `json:"c1-in"`
	C1Jp    string `json:"c1-jp"`
	C1Ru    string `json:"c1-ru"`
	C1Sg    string `json:"c1-sg"`
	C1Uk    string `json:"c1-uk"`
	C1Usc   string `json:"c1-usc"`
	C10De   string `json:"c10-de"`
	C10Eu   string `json:"c10-eu"`
	C10In   string `json:"c10-in"`
	C10Usc  string `json:"c10-usc"`
	C10Use  string `json:"c10-use"`
	C11De   string `json:"c11-de"`
	C11Eu   string `json:"c11-eu"`
	C11In   string `json:"c11-in"`
	C11Usc  string `json:"c11-usc"`
	C12De   string `json:"c12-de"`
	C12Eu   string `json:"c12-eu"`
	C12Use  string `json:"c12-use"`
	C13De   string `json:"c13-de"`
	C13Eu   string `json:"c13-eu"`
	C13Use  string `json:"c13-use"`
	C14De   string `json:"c14-de"`
	C14Use  string `json:"c14-use"`
	C14Usw  string `json:"c14-usw"`
	C15De   string `json:"c15-de"`
	C15Usw  string `json:"c15-usw"`
	C16De   string `json:"c16-de"`
	C16Usw  string `json:"c16-usw"`
	C17De   string `json:"c17-de"`
	C17Usw  string `json:"c17-usw"`
	C18De   string `json:"c18-de"`
	C18Use  string `json:"c18-use"`
	C18Usw  string `json:"c18-usw"`
	C19Use  string `json:"c19-use"`
	C19Usw  string `json:"c19-usw"`
	C2Au    string `json:"c2-au"`
	C2Ca    string `json:"c2-ca"`
	C2De    string `json:"c2-de"`
	C2Eu    string `json:"c2-eu"`
	C2Jp    string `json:"c2-jp"`
	C2Sg    string `json:"c2-sg"`
	C2Uk    string `json:"c2-uk"`
	C2Usc   string `json:"c2-usc"`
	C2Usw   string `json:"c2-usw"`
	C2Za    string `json:"c2-za"`
	C20Use  string `json:"c20-use"`
	C20Usw  string `json:"c20-usw"`
	C21Se   string `json:"c21-use"`
	C21Usw  string `json:"c21-usw"`
	C22Eu   string `json:"c22-eu"`
	C22Use  string `json:"c22-use"`
	C23Eu   string `json:"c23-eu"`
	C23Use  string `json:"c23-use"`
	C24Eu   string `json:"c24-eu"`
	C25Eu   string `json:"c25-eu"`
	C25Use  string `json:"c25-use"`
	C26Eu   string `json:"c26-eu"`
	C26Use  string `json:"c26-use"`
	C27Eu   string `json:"c27-eu"`
	C27Use  string `json:"c27-use"`
	C28Eu   string `json:"c28-eu"`
	C28Use  string `json:"c28-use"`
	C29Eu   string `json:"c29-eu"`
	C29Use  string `json:"c29-use"`
	C3Au    string `json:"c3-au"`
	C3Br    string `json:"c3-br"`
	C3De    string `json:"c3-de"`
	C3Jp    string `json:"c3-jp"`
	C3Sg    string `json:"c3-sg"`
	C3Uk    string `json:"c3-uk"`
	C3Usw   string `json:"c3-usw"`
	C30Eu   string `json:"c30-eu"`
	C30Use  string `json:"c30-use"`
	C31Eu   string `json:"c31-eu"`
	C4De    string `json:"c4-de"`
	C4Eu    string `json:"c4-eu"`
	C4Hk    string `json:"c4-hk"`
	C4Jp    string `json:"c4-jp"`
	C4Uk    string `json:"c4-uk"`
	C4Usc   string `json:"c4-usc"`
	C44Test string `json:"c44-test"`
	C5Eu    string `json:"c5-eu"`
	C5Uk    string `json:"c5-uk"`
	C5Use   string `json:"c5-use"`
	C5Usw   string `json:"c5-usw"`
	C6Ca    string `json:"c6-ca"`
	C6Eu    string `json:"c6-eu"`
	C6Usc   string `json:"c6-usc"`
	C6Use   string `json:"c6-use"`
	C7Eu    string `json:"c7-eu"`
	C7In    string `json:"c7-in"`
	C7Usc   string `json:"c7-usc"`
	C7Use   string `json:"c7-use"`
	C7Usw   string `json:"c7-usw"`
	C8Au    string `json:"c8-au"`
	C8Ca    string `json:"c8-ca"`
	C8Eu    string `json:"c8-eu"`
	C8Usc   string `json:"c8-usc"`
	C8Use   string `json:"c8-use"`
	C8Usw   string `json:"c8-usw"`
	C9Au    string `json:"c9-au"`
	C9Eu    string `json:"c9-eu"`
	C9Usc   string `json:"c9-usc"`
	C9Use   string `json:"c9-use"`
	C9Usw   string `json:"c9-usw"`
	S1Ca    string `json:"s1-ca"`
	S1Sg    string `json:"s1-sg"`
	S1Uae   string `json:"s1-uae"`
	S1Uk    string `json:"s1-uk"`
	S1Usc   string `json:"s1-usc"`
	S2Ca    string `json:"s2-ca"`
	S2De    string `json:"s2-de"`
	S2Uae   string `json:"s2-uae"`
	S2Uk    string `json:"s2-uk"`
	S2Usc   string `json:"s2-usc"`
	S2Use   string `json:"s2-use"`
	S2Usw   string `json:"s2-usw"`
	S2Za    string `json:"s2-za"`
	S3De    string `json:"s3-de"`
	S3Hk    string `json:"s3-hk"`
	S3Jp    string `json:"s3-jp"`
	S3Sg    string `json:"s3-sg"`
	S3Use   string `json:"s3-use"`
	S4Au    string `json:"s4-au"`
	S4Br    string `json:"s4-br"`
	S4Eu    string `json:"s4-eu"`
	S4Hk    string `json:"s4-hk"`
	S4In    string `json:"s4-in"`
	S4Jp    string `json:"s4-jp"`
	S4Sg    string `json:"s4-sg"`
	S5Br    string `json:"s5-br"`
	S5Eu    string `json:"s5-eu"`
	S5In    string `json:"s5-in"`
}

// NewStatusCmd creates and returns a status command to display Algolia server status
func NewStatusCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &StatusOptions{
		IO:     f.IOStreams,
		Config: f.Config,
	}

	cmd := &cobra.Command{
		Use:   "status --incidents",
		Short: "Display Algolia API status",
		Long: heredoc.Doc(`
			This command displays Algolia API status.
		`),
		Example: heredoc.Doc(`
			# Display status
			$ algolia status

			# Display incidents
			$ algolia status --incidents

			# Display status and incidents
			$ algolia status --all

			# Display only specific server(s)
			$ algolia status --server usc,c-10-eu
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			resp, err := http.Get("https://status.algolia.com/1/status")
			if err != nil {
				return fmt.Errorf("An error occured when fetching https://status.algolia.com:", err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("An error occured when reading the body response:", err)
			}

			json.Unmarshal([]byte(body), &opts.Status)

			return runStatusCommand(opts)
		},
	}

	return cmd
}

func runStatusCommand(opts *StatusOptions) error {
	values := reflect.ValueOf(opts.Status.Status)
	typesOf := values.Type()

	for i := 0; i < values.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %v\n", typesOf.Field(i).Name, values.Field(i).Interface())
	}

	return nil
}
