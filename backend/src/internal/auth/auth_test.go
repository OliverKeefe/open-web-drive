package auth

import (
	"context"
	"testing"
)

func TestUserIDFromCtx(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		expectedUserID string
		found          bool
	}{
		{
			name:           "success, user ID exists as string in context",
			ctx:            context.WithValue(context.Background(), userIDKey, "username"),
			expectedUserID: "username",
			found:          true,
		},
		{
			name:           "failure, key missing from context",
			ctx:            context.Background(),
			expectedUserID: "",
			found:          false,
		},
		{
			name:           "failure, key stored as incorrect type in context",
			ctx:            context.WithValue(context.Background(), userIDKey, 11111111),
			expectedUserID: "",
			found:          false,
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			id, ok := UserIDFromCtx(testcase.ctx)
			if id != testcase.expectedUserID {
				t.Errorf("got id %q, want %q", id, testcase.expectedUserID)
			}
			if ok != testcase.found {
				t.Errorf("got ok %t, want %t", ok, testcase.found)
			}
		})
	}
}
