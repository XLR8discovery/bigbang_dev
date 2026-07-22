// Copyright (C) 2026 XLR8discovery PBC
// See LICENSE for copying information.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs/v2"

	"xlr8d.io/bigbang-up/pkg/recipe"
	"xlr8d.io/bigbang-up/pkg/runtime/compose"
	"xlr8d.io/bigbang-up/pkg/runtime/runtime"
	"xlr8d.io/bigbang-up/pkg/runtime/standalone"
)

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "init [<selector>...] OR init <compose|shell> [<selector>...]",
		Args: cobra.MinimumNArgs(1),
		Short: "Initialize new bigbang-up stack with the chosen container orchestrator. " + SelectorHelp + ". Without argument it generates " +
			"full Bigbang cluster with databases (db,minimal,edge)",
	}

	{
		composeCmd := &cobra.Command{
			Use:  "compose [<selector>...]",
			Args: cobra.MinimumNArgs(0),
		}
		composeCmd.RunE = func(cmd *cobra.Command, selector []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			n, err := compose.NewCompose(pwd)
			if err != nil {
				return err
			}
			st, err := recipe.GetStack()
			if err != nil {
				return err
			}
			err = runtime.ApplyRecipes(st, n, normalizedArgs(selector), 0)
			if err != nil {
				return err
			}

			return n.Write()
		}
		cmd.AddCommand(composeCmd)
		cmd.RunE = composeCmd.RunE
	}

	{
		shellCmd := &cobra.Command{
			Use:     "shell [<selector>...]",
			Args:    cobra.MinimumNArgs(0),
			Aliases: []string{"standalone"},
		}
		bigbangProjDir := shellCmd.Flags().StringP("bigbangdir", "s", "", "Directory of the bigbang code.")
		gatewayProjDir := shellCmd.Flags().StringP("gatewaydir", "g", "", "Directory of the gateway code.")
		shellCmd.RunE = func(cmd *cobra.Command, selector []string) error {
			pwd, err := os.Getwd()
			if err != nil {
				return err
			}
			bigbangProjectDir := os.Getenv("BIGBANG_PROJECT_DIR")
			if *bigbangProjDir != "" {
				bigbangProjectDir = *bigbangProjDir
			}
			if bigbangProjectDir == "" {
				return errs.Errorf("Please set \"BIGBANG_PROJECT_DIR\" environment variable or add -s flag with the location of your checked out bigbang/bigbang project. (Required to use web resources")
			}
			gatewayProjectDir := os.Getenv("GATEWAY_PROJECT_DIR")
			if *gatewayProjDir != "" {
				gatewayProjectDir = *gatewayProjDir
			}
			if gatewayProjectDir == "" {
				fmt.Println("WARNING: \"GATEWAY_PROJECT_DIR\" environment variable not set! Please set or add -g flag with the location of your checked out bigbang/gateway-mt project to use web resources.")
				gatewayProjectDir = "/tmp"
			}
			n, err := standalone.NewStandalone(standalone.Paths{
				ScriptDir:  pwd,
				BigbangDir:   bigbangProjectDir,
				GatewayDir: gatewayProjectDir,
				CleanDir:   true,
			})
			if err != nil {
				return err
			}
			st, err := recipe.GetStack()
			if err != nil {
				return err
			}
			err = runtime.ApplyRecipes(st, n, normalizedArgs(selector), 0)
			if err != nil {
				return err
			}

			return n.Write()
		}
		cmd.AddCommand(shellCmd)
	}

	return cmd
}

func normalizedArgs(args []string) []string {
	var res []string
	for _, a := range args {
		for p := range strings.SplitSeq(a, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				res = append(res, p)
			}
		}
	}
	if len(res) == 0 {
		return []string{"db", "minimal", "edge"}
	}
	return res
}

func init() {
	RootCmd.AddCommand(initCmd())
}