package providers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/algolia/cli/pkg/cmdutil"
)

const maxAPIKeyBytes = 64 << 10

func flagSimpleProvider(providerName string) bool {
	switch strings.TrimSpace(providerName) {
	case "openai", "anthropic", "google_genai", "deepseek":
		return true
	default:
		return false
	}
}

func resolveProvidedAPIKey(stdin io.Reader, apiKey, apiKeyEnv string, apiKeyStdin bool) (string, error) {
	n := 0
	if apiKey != "" {
		n++
	}
	if apiKeyStdin {
		n++
	}
	if apiKeyEnv != "" {
		n++
	}
	if n == 0 {
		return "", cmdutil.FlagErrorf(
			"one of --api-key, --api-key-stdin, or --api-key-env is required when creating a provider without -F",
		)
	}
	if n > 1 {
		return "", cmdutil.FlagErrorf("use only one of --api-key, --api-key-stdin, or --api-key-env")
	}

	switch {
	case apiKey != "":
		return apiKey, nil
	case apiKeyStdin:
		b, err := io.ReadAll(io.LimitReader(stdin, maxAPIKeyBytes))
		if err != nil {
			return "", fmt.Errorf("read API key from stdin: %w", err)
		}
		out := strings.TrimSpace(string(b))
		if out == "" {
			return "", cmdutil.FlagErrorf("stdin did not provide a non-empty API key")
		}
		return out, nil
	default:
		raw, ok := os.LookupEnv(apiKeyEnv)
		if !ok {
			return "", cmdutil.FlagErrorf("environment variable %q is not set", apiKeyEnv)
		}
		out := strings.TrimSpace(raw)
		if out == "" {
			return "", cmdutil.FlagErrorf("environment variable %q is empty", apiKeyEnv)
		}
		return out, nil
	}
}

func resolveOptionalAPIKey(stdin io.Reader, apiKey, apiKeyEnv string, apiKeyStdin bool) (string, bool, error) {
	if !apiKeyStdin && apiKeyEnv == "" && apiKey == "" {
		return "", false, nil
	}
	k, err := resolveProvidedAPIKey(stdin, apiKey, apiKeyEnv, apiKeyStdin)
	if err != nil {
		return "", false, err
	}
	return k, true, nil
}

func marshalSimpleProviderCreate(name, providerName, baseURL, apiKey string) ([]byte, error) {
	pn := strings.TrimSpace(providerName)
	if pn == "" {
		return nil, cmdutil.FlagErrorf("--provider must not be empty")
	}
	if !flagSimpleProvider(pn) {
		return nil, cmdutil.FlagErrorf(
			"providers create with flags supports openai, anthropic, google_genai, and deepseek; for %q use -F",
			pn,
		)
	}

	inp := map[string]string{"apiKey": apiKey}
	if baseURL != "" {
		if pn != "openai" && pn != "anthropic" {
			return nil, cmdutil.FlagErrorf("--base-url is only valid with --provider openai or anthropic")
		}
		inp["baseUrl"] = baseURL
	}

	body := struct {
		Name         string            `json:"name"`
		ProviderName string            `json:"providerName"`
		Input        map[string]string `json:"input"`
	}{
		Name:         strings.TrimSpace(name),
		ProviderName: pn,
		Input:        inp,
	}
	return json.Marshal(body)
}

func marshalSimpleProviderPatch(name, baseURL, apiKey string, setKey bool) ([]byte, error) {
	patch := map[string]any{}
	if strings.TrimSpace(name) != "" {
		patch["name"] = strings.TrimSpace(name)
	}
	if setKey || strings.TrimSpace(baseURL) != "" {
		in := map[string]string{}
		if setKey {
			in["apiKey"] = apiKey
		}
		if strings.TrimSpace(baseURL) != "" {
			in["baseUrl"] = strings.TrimSpace(baseURL)
		}
		patch["input"] = in
	}
	if len(patch) == 0 {
		return nil, cmdutil.FlagErrorf(
			"when -F is omitted, specify at least one of --name, --api-key, --api-key-stdin, --api-key-env, or --base-url",
		)
	}
	return json.Marshal(patch)
}

func createShortcutAttempted(name, provider, apiKey, apiKeyEnv, baseURL string, apiKeyStdin bool) bool {
	return strings.TrimSpace(name) != "" ||
		strings.TrimSpace(provider) != "" ||
		apiKey != "" ||
		apiKeyStdin ||
		apiKeyEnv != "" ||
		strings.TrimSpace(baseURL) != ""
}

func inlineFlagsConflictWithFile(file, name, provider, apiKey, apiKeyEnv, baseURL string, apiKeyStdin bool) bool {
	if strings.TrimSpace(file) == "" {
		return false
	}
	return strings.TrimSpace(name) != "" ||
		strings.TrimSpace(provider) != "" ||
		apiKey != "" ||
		apiKeyStdin ||
		apiKeyEnv != "" ||
		strings.TrimSpace(baseURL) != ""
}

func updateUsesInlineFlags(name, apiKey, apiKeyEnv, baseURL string, apiKeyStdin bool) bool {
	return strings.TrimSpace(name) != "" ||
		apiKey != "" ||
		apiKeyStdin ||
		apiKeyEnv != "" ||
		strings.TrimSpace(baseURL) != ""
}

func updateInlineFlagsConflictWithFile(file, name, apiKey, apiKeyEnv, baseURL string, apiKeyStdin bool) bool {
	if strings.TrimSpace(file) == "" {
		return false
	}
	return updateUsesInlineFlags(name, apiKey, apiKeyEnv, baseURL, apiKeyStdin)
}
