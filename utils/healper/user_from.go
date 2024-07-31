package healper

import (
	"context"
	"fmt"
)

const (
	ContextKeyUser    = "USER"
	ContextKeyUserID  = "USER_ID"
	ContextKeySession = "SESSION"
	CreatedBySystem   = "SYSTEM"
)

func GetUserIDFromCtx(ctx context.Context) (userID int64, ok bool) {
	userID, ok = ctx.Value(ContextKeyUserID).(int64)
	return
}

func GetCreatedByFromCtx(ctx context.Context) string {
	userID, ok := GetUserIDFromCtx(ctx)
	if ok {
		return fmt.Sprintf("User %d", userID)
	}

	user, ok := GetUserFromCtx(ctx)
	if ok {
		return fmt.Sprintf("%s", user)
	}

	return CreatedBySystem
}

func GetUserFromCtx(ctx context.Context) (user string, ok bool) {
	user, ok = ctx.Value(ContextKeyUser).(string)
	return
}
