package repositories

import "errors"

var ErrEmailConflict = errors.New("email занят другим пользователем")
