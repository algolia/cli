// Package agentstudio is a thin Go client for Algolia's Agent Studio API
// (https://github.com/algolia/conversational-ai). It uses the same
// X-Algolia-Application-Id / X-Algolia-API-Key auth as the Search API, so it
// plugs into the CLI's existing identity stack without OAuth bearer tokens.
package agentstudio

import (
	"errors"
	"fmt"
	"strings"
)

// Environments accepted by ResolveHost.
const (
	EnvProd    = "prod"
	EnvStaging = "staging"
)

// Regions accepted by ResolveHost.
const (
	RegionEU = "eu"
	RegionUS = "us"
)

// HostOptions controls how the Agent Studio base URL is resolved.
//
// Precedence is: Override > {Region}.algolia.com (per Env) > cluster-proxy
// fallback via ApplicationID. See ResolveHost for details.
type HostOptions struct {
	// Region is the app's hosting region, "eu" or "us". Case-insensitive.
	// Required unless Override or ApplicationID is set.
	Region string

	// Env selects the deployment: "prod" (default) or "staging". Staging is
	// only available in the EU region.
	Env string

	// ApplicationID is used as the last-resort fallback host
	// (https://{appID}.algolia.net/agent-studio), per the Agent Studio README.
	// Useful for legacy profiles created before Region was tracked.
	ApplicationID string

	// Override forces a specific base URL (without trailing slash, without
	// the /1 path suffix). Sourced from --agent-studio-url or
	// ALGOLIA_AGENT_STUDIO_URL. Skips all other resolution.
	Override string
}

// ErrUnknownRegion is returned when Region is set to a value other than
// "eu" or "us" (and no Override or ApplicationID is available).
var ErrUnknownRegion = errors.New("unknown agent studio region")

// ErrStagingNotInRegion is returned when Env=staging is requested for a
// region other than EU. The staging deployment is EU-only.
var ErrStagingNotInRegion = errors.New("agent studio staging is only available in eu")

// ErrNoHostResolvable is returned when none of Override, Region, or
// ApplicationID provides enough information to construct a base URL.
var ErrNoHostResolvable = errors.New(
	"cannot resolve agent studio host: set --agent-studio-url, configure a region on the profile, or pass an application id",
)

// ResolveHost returns the Agent Studio base URL for the given options.
//
// The returned URL has no trailing slash and no /1 suffix; callers append
// "/1/agents" etc. themselves.
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

	// Cluster-proxy fallback: lets the app's own cluster route to the right
	// region. Documented in the Agent Studio README as the recommended URL
	// pattern when a direct regional host isn't known.
	if appID := strings.TrimSpace(opts.ApplicationID); appID != "" {
		return "https://" + appID + ".algolia.net/agent-studio", nil
	}

	return "", ErrNoHostResolvable
}
