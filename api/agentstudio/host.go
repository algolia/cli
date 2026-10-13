// Package agentstudio is a thin Go client for Algolia's Agent Studio
// API (github.com/algolia/conversational-ai). Auth is the standard
// X-Algolia-Application-Id / X-Algolia-API-Key pair — same identity
// stack as the Search API, no OAuth bearer tokens.
package agentstudio

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// DefaultBaseURL is the build-time default for the Agent Studio base
// URL, set via ldflags by `task build`. Empty in production builds
// (cluster-proxy fallback applies); set to the EU staging host for
// internal beta builds. Runtime overrides win.
var DefaultBaseURL string

// EnvAllowInsecureAgentStudioHTTP must be non-empty to permit http://
// overrides (ALGOLIA_AGENT_STUDIO_URL / profile agent_studio_url). Use
// only for local development; production overrides must be https://.
const EnvAllowInsecureAgentStudioHTTP = "ALGOLIA_AGENT_STUDIO_ALLOW_INSECURE_HTTP"

const (
	EnvProd    = "prod"
	EnvStaging = "staging"
)

const (
	RegionEU = "eu"
	RegionUS = "us"
)

// HostOptions controls how the Agent Studio base URL is resolved.
// Precedence: Override > {Region}.algolia.com (per Env) > cluster-proxy
// fallback via ApplicationID.
type HostOptions struct {
	Region        string
	Env           string
	ApplicationID string
	Override      string
}

var (
	// clusterProxyApplicationIDRx construes app IDs that are safe to
	// embed as a single DNS label in https://<id>.algolia.net/agent-studio.
	clusterProxyApplicationIDRx = regexp.MustCompile(`^[A-Za-z0-9]{4,32}$`)

	ErrUnknownRegion      = errors.New("unknown agent studio region")
	ErrStagingNotInRegion = errors.New("agent studio staging is only available in eu")
	ErrNoHostResolvable   = errors.New(
		"cannot resolve agent studio host: set --agent-studio-url, configure a region on the profile, or pass an application id",
	)
)

// ResolveHost returns the Agent Studio base URL for the given options
// (no trailing slash, no /1 suffix — callers append the path).
func ResolveHost(opts HostOptions) (string, error) {
	if opts.Override != "" {
		return normalizeAgentStudioOverride(opts.Override)
	}

	env := strings.ToLower(strings.TrimSpace(opts.Env))
	if env == "" {
		env = EnvProd
	}
	if env != EnvProd && env != EnvStaging {
		return "", fmt.Errorf("unknown agent studio env %q (expected %q or %q)", opts.Env, EnvProd, EnvStaging)
	}

	region := strings.ToLower(strings.TrimSpace(opts.Region))
	switch region {
	case RegionEU:
		if env == EnvStaging {
			return "https://agent-studio.staging.eu.algolia.com", nil
		}
		return "https://agent-studio.eu.algolia.com", nil
	case RegionUS:
		if env == EnvStaging {
			return "", ErrStagingNotInRegion
		}
		return "https://agent-studio.us.algolia.com", nil
	case "":
		// Fall through to cluster-proxy fallback.
	default:
		return "", fmt.Errorf("%w: %q", ErrUnknownRegion, opts.Region)
	}

	// Cluster-proxy fallback: app's own cluster routes to the right region.
	appID := strings.TrimSpace(opts.ApplicationID)
	if appID != "" {
		if err := validateClusterProxyApplicationID(appID); err != nil {
			return "", err
		}
		return "https://" + appID + ".algolia.net/agent-studio", nil
	}

	return "", ErrNoHostResolvable
}

func normalizeAgentStudioOverride(raw string) (string, error) {
	s := strings.TrimRight(strings.TrimSpace(raw), "/")
	if s == "" {
		return "", fmt.Errorf("agent studio url override is empty")
	}
	u, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("agent studio url override: %w", err)
	}
	if u.Scheme == "" {
		return "", fmt.Errorf(
			"agent studio url override must include a scheme (e.g. https://)",
		)
	}
	if u.Host == "" {
		return "", fmt.Errorf("agent studio url override must include a host")
	}
	switch u.Scheme {
	case "https":
		return s, nil
	case "http":
		if os.Getenv(EnvAllowInsecureAgentStudioHTTP) == "" {
			return "", fmt.Errorf(
				"agent studio url must use https:// (got http://); for local development set %s=1",
				EnvAllowInsecureAgentStudioHTTP,
			)
		}
		return s, nil
	default:
		return "", fmt.Errorf(
			"agent studio url scheme %q is not supported (use https://)",
			u.Scheme,
		)
	}
}

func validateClusterProxyApplicationID(appID string) error {
	if !clusterProxyApplicationIDRx.MatchString(appID) {
		return fmt.Errorf(
			"invalid application id %q for agent studio cluster URL: expect 4-32 alphanumeric characters (A-Z, a-z, 0-9)",
			appID,
		)
	}
	return nil
}
