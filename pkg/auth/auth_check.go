package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/config"
	"github.com/algolia/cli/pkg/utils"
)

var WriteAPIKeyDefaultACLs = []string{
	"search",
	"browse",
	"seeUnretrievableAttributes",
	"listIndexes",
	"analytics",
	"logs",
	"addObject",
	"deleteObject",
	"deleteIndex",
	"settings",
	"editSettings",
	"recommendation",
}

// errMissingACLs return an error with the missing ACLs
func errMissingACLs(missing []string) error {
	err := fmt.Sprintf("Missing API Key ACL(s): %s\n", strings.Join(missing, ", "))
	err += "Either edit your profile or use the `--api-key` flag to provide an API Key with the missing ACLs.\n"
	err += "See https://www.algolia.com/doc/guides/security/api-keys/#rights-and-restrictions for more information.\n"
	return errors.New(err)
}

// errAdminAPIKeyRequired is returned when the command requires an admin API Key
var errAdminAPIKeyRequired = errors.New(
	"This command requires an admin API Key. Please use the `--api-key` flag to provide a valid admin API Key.\n",
)

func DisableAuthCheck(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}

	cmd.Annotations["skipAuthCheck"] = "true"
}

func CheckAuth(cfg config.Config) error {
	if cfg.Profile().Name == "" {
		cfg.Profile().LoadDefault()
	}

	_, err := cfg.Profile().GetApplicationID()
	if err != nil {
		return err
	}
	_, err = cfg.Profile().GetAPIKey()
	if err != nil {
		return err
	}

	return nil
}

// CheckACLs check if the current profile has the right ACLs to execute the command
func CheckACLs(cmd *cobra.Command, f *cmdutil.Factory) error {
	if cmd.Annotations == nil {
		return nil
	}

	aclsAsString, ok := cmd.Annotations["acls"]
	if !ok {
		return nil
	}
	neededACLs := strings.Split(aclsAsString, ",")

	client, err := f.SearchClient()
	if err != nil {
		return err
	}
	_, err = client.ListApiKeys()
	if err == nil {
		return nil // Admin API Key, no need to check ACLs
	}

	// Command requires an admin API Key
	if utils.Contains(neededACLs, "admin") {
		return errAdminAPIKeyRequired
	}

	// Check the ACLs of the provided API Key
	keyToUse, err := f.Config.Profile().GetAPIKey()
	if err != nil {
		return err
	}
	apiKey, err := client.GetApiKey(client.NewApiGetApiKeyRequest(keyToUse))
	if err != nil {
		return err
	}

	var acl []string
	for _, a := range apiKey.Acl {
		acl = append(acl, string(a))
	}

	missingACLs := utils.Differences(neededACLs, acl)
	if len(missingACLs) > 0 {
		return errMissingACLs(missingACLs)
	}

	return nil
}

func IsAuthCheckEnabled(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "help", "powershell", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
		return false
	}

	if cmd.Parent() != nil && cmd.Parent().Name() == "completion" {
		return false
	}

	for c := cmd; c.Parent() != nil; c = c.Parent() {
		if c.Annotations != nil && c.Annotations["skipAuthCheck"] == "true" {
			return false
		}
	}

	return true
}
