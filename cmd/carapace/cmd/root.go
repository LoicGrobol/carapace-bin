package cmd

//go:generate go run ../../generate/gen.go

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	spec "github.com/rsteube/carapace-spec"
	"github.com/rsteube/carapace/pkg/ps"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "carapace [flags] [COMPLETER] [bash|elvish|fish|nushell|oil|powershell|tcsh|xonsh|zsh]",
	Long: "multi-shell multi-command argument completer",
	Example: fmt.Sprintf(`  All completers and specs:
    bash:       source <(carapace _carapace bash)
    elvish:     eval (carapace _carapace elvish | slurp)
    fish:       carapace _carapace fish | source
    nushell:    carapace _carapace | save carapace.nu ; nu -c 'source carapace.nu'
    oil:        source <(carapace _carapace oil)
    powershell: carapace _carapace powershell | Out-String | Invoke-Expression
    tcsh:       eval `+"`"+`carapace _carapace tcsh`+"`"+`
    xonsh:      exec($(carapace _carapace xonsh))
    zsh:        source <(carapace _carapace zsh)

  Single completer:
    bash:       source <(carapace chmod bash)
    elvish:     eval (carapace chmod elvish | slurp)
    fish:       carapace chmod fish | source
    nushell:    carapace chmod | save chmod.nu ; nu -c 'source chmod.nu'
    oil:        source <(carapace chmod oil)
    powershell: carapace chmod powershell | Out-String | Invoke-Expression
    tcsh:       eval `+"`"+`carapace _chmod tcsh`+"`"+`
    xonsh:      exec($(carapace chmod xonsh))
    zsh:        source <(carapace chmod zsh)

  Bridge completion:
    bash:       source <(carapace --bridge vault/posener)
    elvish:     eval (carapace --bridge vault/posener|slurp)
    fish:       carapace --bridge vault/posener | source
    nushell:    carapace --bridge vault/posener | save vault.nu ; nu -c 'source vault.nu'
    oil:        source <(carapace --bridge vault/posener)
    powershell: carapace --bridge vault/posener | Out-String | Invoke-Expression
    tcsh:       eval `+"`"+`carapace --bridge vault/posener`+"`"+`
    xonsh:      exec($(carapace --bridge vault/posener))
    zsh:        source <(carapace --bridge vault/posener)
  
  Spec completion:
    bash:       source <(carapace --spec example.yaml)
    elvish:     eval (carapace --spec example.yaml|slurp)
    fish:       carapace --spec example.yaml | source
    oil:        source <(carapace --spec example.yaml)
    nushell:    carapace --spec example.yaml | save example.nu ; nu -c 'source example.nu'
    powershell: carapace --spec example.yaml | Out-String | Invoke-Expression
    tcsh:       eval `+"`"+`carapace --spec example.yaml`+"`"+`
    xonsh:      exec($(carapace --spec example.yaml))
    zsh:        source <(carapace --spec example.yaml)

  Style:
    set:        carapace --style 'carapace.Value=bold,magenta'
    clear:      carapace --style 'carapace.Description='

  Shell parameter is optional and if left out carapace will try to detect it by parent process name.
  Some completions are cached at [%v/carapace].
  Config is written to [%v/carapace].
  Specs are loaded from [%v/carapace/specs].
  `, suppressErr(os.UserCacheDir), suppressErr(os.UserConfigDir), suppressErr(os.UserConfigDir)),
	Args:      cobra.MinimumNArgs(1),
	ValidArgs: completers,
	Run: func(cmd *cobra.Command, args []string) {
		// since flag parsing is disabled do this manually
		switch args[0] {
		case "--bridge":
			if len(args) > 1 {
				// TODO support multiple (comma separated)
				if splitted := strings.SplitN(args[1], "/", 2); len(splitted) == 2 {
					bridgeCompletion(splitted[0], splitted[1], args[2:]...)
				}
			}
		case "--spec":
			if len(args) > 1 {
				specCompletion(args[1], args[2:]...)
			}
		case "--macros":
			if len(args) > 1 {
				printMacro(args[1])
			} else {
				printMacros()
			}
		case "-h":
			cmd.Help()
		case "--help":
			cmd.Help()
		case "-v":
			println(cmd.Version)
		case "--version":
			println(cmd.Version)
		case "--list":
			printCompleters()
		case "--style":
			if len(args) > 1 {
				if err := setStyle(args[1]); err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
				}
			}
		case "_carapace":
			shell := ps.DetermineShell()
			if len(args) > 1 {
				shell = args[1]
			}
			switch shell {
			case "bash":
				fmt.Println(bash_lazy(completers))
			case "bash-ble":
				fmt.Println(bash_ble_lazy(completers))
			case "elvish":
				fmt.Println(elvish_lazy(completers))
			case "fish":
				fmt.Println(fish_lazy(completers))
			case "nushell":
				fmt.Println(nushell_lazy(completers))
			case "oil":
				fmt.Println(oil_lazy(completers))
			case "powershell":
				fmt.Println(powershell_lazy(completers))
			case "tcsh":
				fmt.Println(tcsh_lazy(completers))
			case "xonsh":
				fmt.Println(xonsh_lazy(completers))
			case "zsh":
				fmt.Println(zsh_lazy(completers))
			default:
				fmt.Fprintln(os.Stderr, "could not determine shell")
			}
		default:
			invokeCompleter(args[0])
		}

	},
	FParseErrWhitelist: cobra.FParseErrWhitelist{
		UnknownFlags: true,
	},
	DisableFlagParsing: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func suppressErr(f func() (string, error)) string { s, _ := f(); return s }

func printCompleters() {
	maxlen := 0
	for _, name := range completers {
		if len := len(name); len > maxlen {
			maxlen = len
		}
	}

	for _, name := range completers {
		fmt.Printf("%-"+strconv.Itoa(maxlen)+"v %v\n", name, descriptions[name])
	}
}

func printMacros() {
	maxlen := 0
	names := make([]string, 0)
	for name := range macros {
		names = append(names, name)
		if len := len(name); len > maxlen {
			maxlen = len
		}
	}

	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%-"+strconv.Itoa(maxlen)+"v %v\n", name, macroDescriptions[name])
	}
}

func printMacro(name string) {
	if m, ok := macros[name]; ok {
		path := strings.Replace(name, ".", "/", -1)
		signature := ""
		if s := m.Signature(); s != "" {
			signature = fmt.Sprintf("(%v)", s)
		}

		fmt.Printf(`signature:   $_%v%v
description: %v
reference:   https://pkg.go.dev/github.com/rsteube/carapace-bin/pkg/actions/%v#Action%v
`, name, signature, macroDescriptions[name], filepath.Dir(path), filepath.Base(path))
	}
}

func invokeCompleter(completer string) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	os.Args[1] = "_carapace"
	executeCompleter(completer)

	w.Close()
	out := <-outC
	os.Stdout = old

	executable, err := os.Executable()
	if err != nil {
		panic(err.Error()) // TODO exit with error message
	}
	executableName := filepath.Base(executable)
	patched := strings.Replace(string(out), fmt.Sprintf("%v _carapace", executableName), fmt.Sprintf("%v %v", executableName, completer), -1)      // general callback
	patched = strings.Replace(patched, fmt.Sprintf("'%v', '_carapace'", executableName), fmt.Sprintf("'%v', '%v'", executableName, completer), -1) // xonsh callback
	fmt.Print(patched)

}

func setStyle(s string) error {
	if splitted := strings.SplitN(s, "=", 2); len(splitted) == 2 {
		return style.Set(splitted[0], splitted[1])
	}
	return fmt.Errorf("invalid format: '%v'", s)
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().String("bridge", "", "bridge completion")
	rootCmd.Flags().Bool("list", false, "list completers")
	rootCmd.Flags().String("spec", "", "spec completion")
	rootCmd.Flags().String("style", "", "set style")

	for m, f := range macros {
		spec.AddMacro(m, f)
	}
}
