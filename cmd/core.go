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

const version = "0.4.0"

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

				theme, _ := cmd.Flags().GetString("theme")
				validateThemeSupport(theme)
				if cmd.Annotations == nil {
					cmd.Annotations = make(map[string]string)
				}
				cmd.Annotations["theme"] = theme
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
			run(program.NewProgram("grep", true), cmd.Annotations["theme"])
		},
	}

	sedCmd = &cobra.Command{
		Use:   "sed",
		Short: `Play with sed`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("sed", true), cmd.Annotations["theme"])
		},
	}

	awkCmd = &cobra.Command{
		Use:   "awk",
		Short: `Play with awk`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("awk", true), cmd.Annotations["theme"])
		},
	}

	jqCmd = &cobra.Command{
		Use:   "jq",
		Short: `Play with jq`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("jq", false), cmd.Annotations["theme"])
		},
	}

	yqCmd = &cobra.Command{
		Use:   "yq",
		Short: `Play with yq`,
		Run: func(cmd *cobra.Command, args []string) {
			run(program.NewProgram("yq", false), cmd.Annotations["theme"])
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

func validateThemeSupport(theme string) {
	var validThemes []string

	_, exists := ui.Themes[theme]
	if !exists {
		for theme := range ui.Themes {
			validThemes = append(validThemes, theme)
		}
		fmt.Printf("Error: Invalid theme '%s'. Valid themes are: %v\n", theme, validThemes)
		os.Exit(1)
	}
}

func validateProgramExists(program string) {
	_, err := exec.LookPath(program)
	if err != nil {
		exitWithError(program + " not found")
	}
}

func run(program program.Program, theme string) error {

	var userInterface *ui.UI
	var stdinTmpFile string
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
				stdin = append(stdin, '\n')
			}
			input = string(stdin)

			tmpFile, err := os.CreateTemp("", "play")
			if err != nil {
				exitWithError(err)
			}
			defer tmpFile.Close()

			_, err = tmpFile.Write([]byte(input))
			if err != nil {
				exitWithError(err)
			}
			stdinTmpFile = tmpFile.Name()
		}
	}

	userInterface = ui.NewUI(program.Name, program.RespectsEndOfOptions, stdinTmpFile, theme)
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
	rootCmd.PersistentFlags().String("theme", "monokai", "theme")
	versionCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		command.Flags().MarkHidden("theme")
		command.Parent().HelpFunc()(command, strings)
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitWithError(err)
	}
}
