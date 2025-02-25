package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace-bin/pkg/actions/tools/golang"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "compile packages and dependencies",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(buildCmd).Standalone()

	buildCmd.Flags().BoolS("a", "a", false, "force rebuilding of packages that are already up-to-date.")
	buildCmd.Flags().String("asmflags", "", "arguments to pass on each go tool asm invocation")
	buildCmd.Flags().String("buildmode", "", "build mode to use")
	buildCmd.Flags().String("compiler", "", "name of compiler to use")
	buildCmd.Flags().String("gccgoflags", "", "arguments to pass on each gccgo compiler/linker invocation")
	buildCmd.Flags().String("gcflags", "", "arguments to pass on each go tool compile invocation.")
	buildCmd.Flags().BoolS("i", "i", false, "install the packages that are dependencies of the target")
	buildCmd.Flags().String("installsuffix", "", "a suffix to use in the name of the package installation directory")
	buildCmd.Flags().String("ldflags", "", "arguments to pass on each go tool link invocation")
	buildCmd.Flags().Bool("linkshared", false, "build code that will be linked against shared libraries")
	buildCmd.Flags().String("mod", "", "module download mode to use")
	buildCmd.Flags().Bool("modcacherw", false, "leave newly-created directories in the module cache read-write")
	buildCmd.Flags().String("modfile", "", "read and possibly write an alternate go.mod file")
	buildCmd.Flags().Bool("msan", false, "enable interoperation with memory sanitizer")
	buildCmd.Flags().BoolS("n", "n", false, "print the commands but do not run them.")
	buildCmd.Flags().StringS("o", "o", "", "set output file or directory")
	buildCmd.Flags().StringS("p", "p", "", "the number of programs to run in parallel")
	buildCmd.Flags().String("pkgdir", "", "install and load all packages from dir")
	buildCmd.Flags().Bool("race", false, "enable data race detection")
	buildCmd.Flags().String("tags", "", "a comma-separated list of build tags to consider satisfied during the")
	buildCmd.Flags().String("toolexec", "", "a program to use to invoke toolchain programs like vet and asm")
	buildCmd.Flags().Bool("trimpath", false, "remove all file system paths from the resulting executable")
	buildCmd.Flags().BoolS("v", "v", false, "print the names of packages as they are compiled")
	buildCmd.Flags().Bool("work", false, "print the name of the temporary work directory")
	buildCmd.Flags().BoolS("x", "x", false, "print the commands.")
	rootCmd.AddCommand(buildCmd)

	carapace.Gen(buildCmd).FlagCompletion(carapace.ActionMap{
		"buildmode": carapace.ActionValues("archive", "c-archive", "c-shared", "default", "shared", "exe", "pie", "plugin"),
		"compiler":  carapace.ActionValues("gccgo", "gc"),
		"mod":       carapace.ActionValues("readonly", "vendor", "mod"),
		"modfile":   carapace.ActionFiles(".mod"),
		"n":         carapace.ActionValues("1", "2", "3", "4", "5", "6", "7", "8"),
		"o":         carapace.ActionFiles(),
		"pkgdir":    carapace.ActionDirectories(),
		"tags": carapace.ActionMultiParts(",", func(c carapace.Context) carapace.Action {
			return golang.ActionBuildTags().Invoke(c).Filter(c.Parts).ToA()
		}),
	})
}
