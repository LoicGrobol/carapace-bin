package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bin/pkg/actions/bridge"
	"github.com/spf13/cobra"
)

var buildxCmd = &cobra.Command{
	Use:                "buildx",
	Short:              "Extended build capabilities with BuildKit",
	Run:                func(cmd *cobra.Command, args []string) {},
	DisableFlagParsing: true,
}

func init() {
	carapace.Gen(buildxCmd).Standalone()

	rootCmd.AddCommand(buildxCmd)

	carapace.Gen(buildxCmd).PositionalAnyCompletion(
		bridge.ActionCarapaceBin("docker-buildx"),
	)
}
