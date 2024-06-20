package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/config"
	"github.com/joeldotdias/weave/internal/tui/multiChoice"
	"github.com/joeldotdias/weave/internal/tui/textArea"
	"github.com/joeldotdias/weave/internal/tui/textInput"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Opts struct {
	add_all     bool
	Title       *textInput.Response
	Symbol      *multiChoice.Selected
	Description *textArea.Description
}

var rootCmd = &cobra.Command{
	Use:   "weave",
	Short: "A tool to write better commit messages",
	Long: `Weave provides an intuituive TUI to write descriptive
commit messages when needed.
Configuration can be written at $XDG_CONFIG_HOME/.config/weave/config.toml`,

	Run: func(cmd *cobra.Command, args []string) {
		checkIfInsideGitRepo()

		conf := config.MakeConfig()
		opts := Opts{
			add_all:     conf.Add_all,
			Title:       &textInput.Response{},
			Symbol:      &multiChoice.Selected{},
			Description: &textArea.Description{},
		}

		var (
			tprogram    *tea.Program
			title, desc string
			exit        bool
			err         error
		)

		if len(conf.Title) != 0 {
			title = conf.Title
		} else {
			tprogram = tea.NewProgram(textInput.InitTextInputModel(opts.Title, "Your title here...", &exit))
			if _, err = tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}
			checkExit(tprogram, exit)

			tprogram = tea.NewProgram(multiChoice.InitMultiChoiceModel(conf.SymbolChoices(conf.Format), opts.Symbol, "Choose symbol", &exit))
			if _, err = tprogram.Run(); err != nil {
				cobra.CheckErr(err)
			}
			checkExit(tprogram, exit)
			title = fmt.Sprintf("%s%s %s", opts.Symbol.Value(), conf.Separator, opts.Title.Value())
		}

		tprogram = tea.NewProgram(textArea.InitTextAreaModel(opts.Description, "Your description here", "Limit this to 72 words", &exit))
		if _, err = tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
		checkExit(tprogram, exit)
		desc = opts.Description.Value()

		if err = makeCommit(title, desc, conf.Add_all); err != nil {
			fmt.Println("Whoops, you might not have added anything to commit")
		}
	},
}

func Execute() {
	bindFlags()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.MakePresets()

	rootCmd.PersistentFlags().BoolP("add_all", "a", false, "Add all files before committing")
	rootCmd.PersistentFlags().StringP("title", "t", "", "Skip the title step by adding your own")
	rootCmd.PersistentFlags().StringP("format", "f", "<type> <symbol>", "Format the title prefix by using <type> and <symbol>")
	rootCmd.PersistentFlags().StringP("separator", "s", ": ", "Separator between the prefix and title")
}

func bindFlags() {
	err := viper.BindPFlag("add_all", rootCmd.PersistentFlags().Lookup("add_all"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("title", rootCmd.PersistentFlags().Lookup("title"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("format", rootCmd.PersistentFlags().Lookup("format"))
	if err != nil {
		log.Fatal(err)
	}
	err = viper.BindPFlag("separator", rootCmd.PersistentFlags().Lookup("separator"))
	if err != nil {
		log.Fatal(err)
	}
}

func checkExit(tprogram *tea.Program, exit bool) {
	if exit {
		if err := tprogram.ReleaseTerminal(); err != nil {
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
		fmt.Scan(&initRepo)

		if initRepo == "y" || initRepo == "Y" {
			_, err = exec.Command("git", "init").Output()
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			fmt.Println("You need to be inside a git repo to make commits")
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

	_, err = exec.Command("git", "commit", "-m", title, "-m", desc).Output()
	return err
}
