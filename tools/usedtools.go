//go:build tools

package tools

import (
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/go-task/task/v3/cmd/task"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/mailru/easyjson/easyjson"
)
