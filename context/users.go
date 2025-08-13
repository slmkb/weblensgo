package context

import (
	"context"

	"github.com/slmkb/weblensgo/models"
)

type key string

const (
	userKey      key = "user"
	galleriesKey key = "galleries"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		return nil
	}
	return user
}

// func WithGalleries(ctx context.Context, galleries *map[string]string) context.Context {
// 	return context.WithValue(ctx, galleriesKey, galleries)
// }

// func Galleries(ctx context.Context) *map[string]string {
// 	val := ctx.Value(galleriesKey)
// 	galleries, ok := val.(*map[string]string)
// 	if !ok {
// 		return nil
// 	}
// 	return galleries
// }
