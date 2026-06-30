// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package modify

import (
	"github.com/spf13/cobra"

	"xlr8d.io/bigbang-up/cmd"
	"xlr8d.io/bigbang-up/pkg/common"
	"xlr8d.io/bigbang-up/pkg/recipe"
	"xlr8d.io/bigbang-up/pkg/runtime/runtime"
)

func imageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "image <selector>... <image>",
		Short: "Change container image for one more more services",
		Long:  "Change image of specified services." + cmd.SelectorHelp,
		Args:  cobra.MinimumNArgs(2),
		RunE:  cmd.ExecuteBigbangUP(setImage),
	}
}

func init() {
	cmd.RootCmd.AddCommand(imageCmd())
}

func setImage(st recipe.Stack, rt runtime.Runtime, args []string) error {
	selector, image := common.SplitArgsSelector1(args)
	return runtime.ModifyService(st, rt, selector, func(s runtime.Service) error {
		return s.ChangeImage(func(s string) string {
			return image
		})
	})
}