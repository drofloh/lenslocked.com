package context

import (
	"context"

	"github.com/drofloh/lenslocked.com/models"
)

const (
	userKey priviateKey = "user"
)

type priviateKey string

// WithUser ...
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User ...
func User(ctx context.Context) *models.User {
	if temp := ctx.Value("user"); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
