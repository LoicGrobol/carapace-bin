package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Delete all folders in the clean-targets list",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(cleanCmd).Standalone()

	cleanCmd.Flags().String("profile", "", "Which profile to load.")
	cleanCmd.Flags().String("profiles-dir", "", "Which directory to look in for the profiles.yml file.")
	cleanCmd.Flags().String("project-dir", "", "Which directory to look in for the dbt_project.yml file")
	cleanCmd.Flags().StringP("target", "t", "", "Which target to load for the given profile")
	cleanCmd.Flags().String("vars", "", "Supply variables to the project")
	rootCmd.AddCommand(cleanCmd)
}
