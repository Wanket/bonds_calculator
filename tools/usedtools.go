//go:build tools

package tools

import (
	_ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/mailru/easyjson/easyjson"
)
