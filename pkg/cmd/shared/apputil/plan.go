package apputil

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/algolia/cli/api/dashboard"
	"github.com/algolia/cli/pkg/cmdutil"
	"github.com/algolia/cli/pkg/prompt"
)

// ResolvePlan maps a --plan value to one of the available plans. An exact match
// on the plan id wins; the user-facing "free" choice maps to the free-type
// template, whose id is not fixed (it can be "build"), so it is matched on type.
func ResolvePlan(plans []dashboard.Plan, value string) (*dashboard.Plan, error) {
	for i := range plans {
		if plans[i].ID == value {
			return &plans[i], nil
		}
	}
	if value == dashboard.PlanTypeFree {
		if free := FindFreePlan(plans); free != nil {
			return free, nil
		}
	}
	return nil, cmdutil.FlagErrorf(
		"invalid plan %q; available plans: %s",
		value,
		strings.Join(PlanChoices(plans), ", "),
	)
}

var knownPaidPlanNames = map[string]string{
	"grow":      "Grow",
	"grow-plus": "Grow Plus",
}

// KnownPaidPlan recognizes a documented paid --plan value even when the
// self-serve endpoint omits it (paid plans appear only once billing is on file).
func KnownPaidPlan(value string) *dashboard.Plan {
	name, ok := knownPaidPlanNames[value]
	if !ok {
		return nil
	}
	return &dashboard.Plan{
		ID:   value,
		Name: name,
		Type: "freeform",
	}
}

// PlanAvailable reports whether a plan with the given id is in the list.
func PlanAvailable(plans []dashboard.Plan, id string) bool {
	for i := range plans {
		if plans[i].ID == id {
			return true
		}
	}
	return false
}

// PlanChoices returns the user-facing plan identifiers (the free plan is shown
// as "free" regardless of its underlying id).
func PlanChoices(plans []dashboard.Plan) []string {
	choices := make([]string, 0, len(plans))
	for _, p := range plans {
		if p.IsFree() {
			choices = append(choices, dashboard.PlanTypeFree)
		} else {
			choices = append(choices, p.ID)
		}
	}
	return choices
}

// PlanTelemetryID returns the plan identifier used in telemetry properties:
// the user-facing "free" for the free tier (whose underlying template id is
// not fixed and can be "build"), the plan id otherwise.
func PlanTelemetryID(p dashboard.Plan) string {
	if p.IsFree() {
		return dashboard.PlanTypeFree
	}
	return p.ID
}

// SelectablePlans returns the plans a user may choose from. When hideNonFree is
// true (no payment method on file) only the free plan(s) are offered, because
// paid plans require billing details the CLI can't collect.
func SelectablePlans(plans []dashboard.Plan, hideNonFree bool) []dashboard.Plan {
	if !hideNonFree {
		return plans
	}
	free := make([]dashboard.Plan, 0, 1)
	for _, p := range plans {
		if p.IsFree() {
			free = append(free, p)
		}
	}
	return free
}

// FindFreePlan returns the free-tier plan, or nil if none is present.
func FindFreePlan(plans []dashboard.Plan) *dashboard.Plan {
	for i := range plans {
		if plans[i].IsFree() {
			return &plans[i]
		}
	}
	return nil
}

// PickPlan shows an interactive selector over the candidate plans.
func PickPlan(candidates []dashboard.Plan) (*dashboard.Plan, error) {
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no plans are available")
	}
	labels := make([]string, len(candidates))
	for i, p := range candidates {
		labels[i] = fmt.Sprintf("%s — %s", p.Name, p.Price)
	}
	var selected int
	if err := prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select a plan:",
			Options: labels,
		},
		&selected,
	); err != nil {
		return nil, err
	}
	return &candidates[selected], nil
}
