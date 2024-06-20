package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/config"
	"github.com/joeldotdias/weave/internal/tui/multiChoice"
	"github.com/joeldotdias/weave/internal/tui/textArea"
	"github.com/joeldotdias/weave/internal/tui/textInput"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "weave",
	Short: "A tool to write better commit messages",
	Long: `Weave provides an intuituive TUI to write descriptive
commit messages when needed.
Configuration can be written at $XDG_CONFIG_HOME/.config/weave/config.toml`,
	Run: weave,
}

func weave(cmd *cobra.Command, args []string) {
	conf := config.MakeConfig()
	opts := Opts{
		Title:       &textInput.Response{},
		Symbol:      &multiChoice.Selected{},
		Description: &textArea.Description{},
	}
	var tprogram *tea.Program

	tprogram = tea.NewProgram(textInput.InitTextInputModel(opts.Title, "Your title here..."))

	var title, desc string
	if len(conf.Title) != 0 {
		title = conf.Title
	} else {
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
		title = opts.Title.Value()
	}

	multiOptions := SymbolsList{
		options: conf.SymbolChoices(),
	}

	tprogram = tea.NewProgram(multiChoice.InitMultiChoiceModel(multiOptions.options, opts.Symbol, "Choose symbol"))
	if _, err := tprogram.Run(); err != nil {
		cobra.CheckErr(err)
	}
	symbol := opts.Symbol.Value()

	tprogram = tea.NewProgram(textArea.InitTextAreaModel(opts.Description, "Your description here", "Limit this to 72 words"))
	if _, err := tprogram.Run(); err != nil {
		cobra.CheckErr(err)
	}
	desc = opts.Description.Value()

	fmt.Printf("\nYour message: %s: %s", symbol, title)
	fmt.Printf("\nYour Description:\n%s\n", desc)
}

func Execute() {
	bindFlags()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type Opts struct {
	Title       *textInput.Response
	Symbol      *multiChoice.Selected
	Description *textArea.Description
}

type SymbolsList struct {
	options []string
}

var (
	add_all bool
	title   string
)

func init() {
	config.MakePresets()

	rootCmd.PersistentFlags().BoolVarP(&add_all, "add_all", "a", false, "Add all files before committing")
	rootCmd.PersistentFlags().StringVarP(&title, "title", "t", "", "Skip the title step by adding one now")
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
}
