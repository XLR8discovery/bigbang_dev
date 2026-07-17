// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeebo/errs/v2"

	"xlr8d.io/bigbang-up/pkg/runtime/compose"
	"xlr8d.io/bigbang-up/pkg/runtime/runtime"
	"xlr8d.io/bigbang-up/pkg/runtime/standalone"
)

// FromDir creates the right runtime based on available file names in the directory.
func FromDir(dir string) (runtime.Runtime, error) {
	_, err := os.Stat(filepath.Join(dir, "docker-compose.yaml"))
	if err == nil {
		return compose.NewCompose(dir)
	}

	_, err = os.Stat(filepath.Join(dir, "supervisord.conf"))
	if err == nil {
		bigbangProjectDir := os.Getenv("BIGBANG_PROJECT_DIR")
		if bigbangProjectDir == "" {
			return nil, errs.Errorf("Please set \"BIGBANG_PROJECT_DIR\" environment variable with the location of your checked out bigbang/bigbang project. (Required to use web resources")
		}
		gatewayProjectDir := os.Getenv("GATEWAY_PROJECT_DIR")
		if gatewayProjectDir == "" {
			fmt.Println("WARNING: \"GATEWAY_PROJECT_DIR\" environment variable not set! Please set or add -g flag with the location of your checked out bigbang/gateway-mt project to use web resources.")
			gatewayProjectDir = "/tmp"
		}
		return standalone.NewStandalone(standalone.Paths{
			ScriptDir:  dir,
			BigbangDir:   bigbangProjectDir,
			GatewayDir: gatewayProjectDir,
			CleanDir:   false,
		})
	}

	return nil, errors.New("directory doesn't contain supported deployment descriptor")
}