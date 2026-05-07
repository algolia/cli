package shared

import "context"

// OrBackground returns ctx, or context.Background() when ctx is nil.
// Cobra always passes a real context, but tests sometimes invoke run
// helpers directly with a nil opts.Ctx; this keeps callers from having
// to special-case that.
func OrBackground(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
