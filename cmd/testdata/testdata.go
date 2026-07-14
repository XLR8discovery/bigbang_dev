// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package testdata

import (
	_ "embed"
	"os"
	"path/filepath"

	"xlr8d.io/bigbang-up/pkg/recipe"
	"xlr8d.io/bigbang-up/pkg/runtime/compose"
	"xlr8d.io/bigbang-up/pkg/runtime/runtime"
)

//go:embed docker-compose.yaml
var composeFile []byte

func InitCompose(dir string) (st recipe.Stack, rt runtime.Runtime, err error) {
	err = os.WriteFile(filepath.Join(dir, "docker-compose.yaml"), composeFile, 0644)
	if err != nil {
		return
	}
	rt, err = compose.NewCompose(dir)
	if err != nil {
		return
	}

	st, err = recipe.GetStack()
	if err != nil {
		return
	}
	err = runtime.ApplyRecipes(st, rt, []string{"db", "minimal"}, 0)
	return
}
