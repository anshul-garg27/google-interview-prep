package leaderboard

import (
	"coffee/constants"
	"coffee/core/appcontext"
	"context"
	"errors"
)

func validateRequest(ctx context.Context, page int, size int) error {
	appCtx := ctx.Value(constants.AppContextKey).(*appcontext.RequestContext)
	if (appCtx.PlanType == nil || *appCtx.PlanType == constants.FreePlan) && (page > 1 || size > 10) {
		err := errors.New("invalid Request")
		return err
	}
	return nil
}
