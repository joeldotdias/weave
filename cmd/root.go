package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/config"
	"github.com/joeldotdias/weave/internal/tui"
	"github.com/joeldotdias/weave/internal/tui/multiChoice"
	"github.com/joeldotdias/weave/internal/tui/textArea"
	"github.com/joeldotdias/weave/internal/tui/textInput"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Opts struct {
	Title       *textInput.Response
	Symbol      *multiChoice.Selected
	Description *textArea.Description
}

var rootCmd = &cobra.Command{
	Use:     "weave",
	Version: "v0.1.0",
	Short:   "A tool to write better commit messages",
	Long: `Weave provides an intuituive TUI to write descriptive
commit messages when needed.
Configuration can be written at $XDG_CONFIG_HOME/weave/config.toml`,

	Run: func(cmd *cobra.Command, args []string) {
		checkIfInsideGitRepo()

		conf := config.MakeConfig()
		opts := Opts{
			Title:       &textInput.Response{},
			Symbol:      &multiChoice.Selected{},
			Description: &textArea.Description{},
		}

		var (
			prog        *tea.Program
			title, desc string
			exit        bool
			err         error
		)

		if len(conf.Title) != 0 {
			title = conf.Title
		} else {
			prog = tea.NewProgram(textInput.InitTextInputModel(opts.Title, "Your title here...", &conf.Theme, &exit))
			if _, err = prog.Run(); err != nil {
				cobra.CheckErr(err)
			}
			checkExit(prog, exit)

			prog = tea.NewProgram(multiChoice.InitMultiChoiceModel(conf.SymbolChoices(conf.Format), opts.Symbol, "Choose your prefix", &conf.Theme, &exit))
			if _, err = prog.Run(); err != nil {
				cobra.CheckErr(err)
			}
			checkExit(prog, exit)
			title = fmt.Sprintf("%s%s %s", opts.Symbol.Value(), conf.Separator, opts.Title.Value())
		}

		prog = tea.NewProgram(textArea.InitTextAreaModel(opts.Description, "Your description here", "Limit this to 72 words", &conf.Theme, &exit))
		if _, err = prog.Run(); err != nil {
			cobra.CheckErr(err)
		}
		checkExit(prog, exit)
		desc = opts.Description.Value()

		fmt.Println(tui.NotifStyle.Render("Making your commit now"))
		if err = makeCommit(title, desc, conf.Add_all); err != nil {
			fmt.Println(tui.ErrStyle.Render("Whoops, you might not have added anything to commit"))
		}
	},
}

func Execute() {
	var err error

	flags := []string{"add_all", "title", "format", "separator"}
	for _, flag := range flags {
		err = viper.BindPFlag(flag, rootCmd.PersistentFlags().Lookup(flag))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s\n" .Version}}`)

	config.MakePresets()

	rootCmd.PersistentFlags().BoolP("add_all", "a", false, "Add all files before committing")
	rootCmd.PersistentFlags().StringP("title", "t", "", "Skip the title step by adding your own")
	rootCmd.PersistentFlags().StringP("format", "f", "<type> <symbol>", "Format the title prefix by using <type> and <symbol>")
	rootCmd.PersistentFlags().StringP("separator", "s", ": ", "Separator between the prefix and title")
}

func checkExit(prog *tea.Program, exit bool) {
	if exit {
		if err := prog.ReleaseTerminal(); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func checkIfInsideGitRepo() {
	var (
		err      error
		exitErr  *exec.ExitError
		initRepo string
	)

	_, err = exec.Command("git", "rev-parse", "--git-dir").Output()
	if err != nil && !errors.As(err, &exitErr) {
		log.Fatalf("%p\nExiting...", err.Error)
	}

	// if an error (except for an ExitError) is rerturned,
	// then we aren't inside a git repo
	if err != nil {
		fmt.Print("You aren't inside a git repo. Initialize one now?(y/n) ")
		_, err = fmt.Scan(&initRepo)
		if err != nil {
			log.Fatal(err)
		}

		if initRepo == "y" || initRepo == "Y" {
			_, err = exec.Command("git", "init").Output()
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			fmt.Println(tui.ErrStyle.Render("You need to be inside a git repo to make commits"))
			os.Exit(1)
		}
	}
}

func makeCommit(title string, desc string, add_all bool) error {
	var err error

	if add_all {
		_, err = exec.Command("git", "add", ".").Output()
		if err != nil {
			return err
		}
	}

	out, err := exec.Command("git", "commit", "-m", title, "-m", desc).Output()
	fmt.Println(string(out))
	return err
}
