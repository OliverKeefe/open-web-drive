package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func UserIDFromCtx(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func ClaimsFromCtx(ctx context.Context) (map[string]interface{}, bool) {
	claims, ok := ctx.Value(claimsKey).(map[string]interface{})
	return claims, ok
}

func HasClaim(ctx context.Context, userID uuid.UUID, requiredClaim string) (bool, error) {
	userIDStoredInCtx, ok := UserIDFromCtx(ctx)
	if !ok {
		return false, errors.New("unauthenticated: no user ID found in context")
	}

	if userIDStoredInCtx != userID.String() {
		return false, errors.New("unauthorized: userID mismatch")
	}

	claims, ok := ClaimsFromCtx(ctx)
	if !ok {
		return false, errors.New("unable to get claims from context")
	}

	scopes, present := claims["scopes"]
	if !present {
		return false, errors.New("no scope claims found in token")
	}

	scopesStr, ok := scopes.(string)
	if !ok {
		return false, errors.New("couldn't convert scope claim to string")
	}

	hasClaim := false
	for _, scope := range strings.Fields(scopesStr) {
		if scope == requiredClaim {
			hasClaim = true
			break
		}
	}

	if !hasClaim {
		return false, errors.New("unauthorized: missing required scope")
	}

	return true, nil
}
