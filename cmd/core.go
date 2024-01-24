package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	ui "github.com/paololazzari/play/src/ui"
	program "github.com/paololazzari/play/src/util"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const version = "0.3.3"

func completionCommand() *cobra.Command {
	return &cobra.Command{
		Use: "completion",
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "play",
		Short: "play",
		Long:  `play`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use != "version" {
				validateProgramExists(cmd.Use)
			}
			validateArgs(args)
			return nil
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Long:  `Print version number`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("play", version)
		},
	}

	grepCmd = &cobra.Command{
		Use:   "grep",
		Short: `Play with grep`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("grep", true))
		},
	}

	sedCmd = &cobra.Command{
		Use:   "sed",
		Short: `Play with sed`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("sed", true))
		},
	}

	awkCmd = &cobra.Command{
		Use:   "awk",
		Short: `Play with awk`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("awk", true))
		},
	}

	jqCmd = &cobra.Command{
		Use:   "jq",
		Short: `Play with jq`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("jq", false))
		},
	}

	yqCmd = &cobra.Command{
		Use:   "yq",
		Short: `Play with yq`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("yq", false))
		},
	}
)

func exitWithError(e interface{}) {
	fmt.Fprintln(os.Stderr, e)
	os.Exit(1)
}

func validateArgs(args []string) {
	if len(args) > 0 {
		exitWithError("Invalid number of arguments")
	}
}

func validateProgramExists(program string) {
	_, err := exec.LookPath(program)
	if err != nil {
		exitWithError(program + " not found")
	}
}

func run(program program.Program) error {

	var userInterface *ui.UI
	input := ""

	// check whether file descriptor is terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		if runtime.GOOS != "windows" {
			_, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
			if err != nil {
				exitWithError(err)
			}
			var stdin []byte
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				stdin = append(stdin, scanner.Bytes()...)
			}
			input = string(stdin)
		}
	}

	userInterface = ui.NewUI(program.Name, program.RespectsEndOfOptions, input)
	userInterface.InitUI()
	userInterface.Run()
	return nil
}

func init() {
	completion := completionCommand()
	completion.Hidden = true
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(grepCmd)
	rootCmd.AddCommand(sedCmd)
	rootCmd.AddCommand(awkCmd)
	rootCmd.AddCommand(jqCmd)
	rootCmd.AddCommand(yqCmd)
	rootCmd.AddCommand(completion) // https://github.com/spf13/cobra/issues/1507
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}
