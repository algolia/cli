// Package agentstudio is a thin Go client for Algolia's Agent Studio
// API (github.com/algolia/conversational-ai). Auth is the standard
// X-Algolia-Application-Id / X-Algolia-API-Key pair — same identity
// stack as the Search API, no OAuth bearer tokens.
package agentstudio

import (
	"errors"
	"fmt"
	"strings"
)

// DefaultBaseURL is the build-time default for the Agent Studio base
// URL, set via ldflags by `task build`. Empty in production builds
// (cluster-proxy fallback applies); set to the EU staging host for
// internal beta builds. Runtime overrides win.
var DefaultBaseURL string

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
		return strings.TrimRight(opts.Override, "/"), nil
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
	if appID := strings.TrimSpace(opts.ApplicationID); appID != "" {
		return "https://" + appID + ".algolia.net/agent-studio", nil
	}

	return "", ErrNoHostResolvable
}
