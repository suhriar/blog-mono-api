package utils

import (
	"context"
	"errors"

	"github.com/suhriar/blog-mono-api/model"
)

func GetUserFromContext(ctx context.Context) (user model.UserAuth, err error) {
	username, ok1 := ctx.Value(model.UserNameKey).(string)
	email, ok2 := ctx.Value(model.UserEmailKey).(string)
	id, ok3 := ctx.Value(model.UserIDlKey).(int64)

	if !ok1 || !ok2 || !ok3 {
		return user, errors.New("could not get user data from context")
	}

	user.ID = id
	user.Username = username
	user.Email = email

	return user, nil
}
