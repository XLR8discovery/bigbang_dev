// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package cmd

import (
	"github.com/spf13/cobra"

	"xlr8d.io/bigbang-up/pkg/recipe"
	"xlr8d.io/bigbang-up/pkg/runtime/runtime"
)

func init() {
	var instance int
	cmd := &cobra.Command{
		Use:   "add <selector>...",
		Args:  cobra.MinimumNArgs(1),
		Short: "add more services to existing stack. " + SelectorHelp,
		RunE: ExecuteBigbangUP(func(stack recipe.Stack, rt runtime.Runtime, selector []string) error {
			return runtime.ApplyRecipes(stack, rt, selector, instance)
		}),
	}

	cmd.PersistentFlags().IntVarP(&instance, "instance", "i", 0, "Number of requested instance (default/0 = use the one defined in the recipe")
	RootCmd.AddCommand(cmd)

}
