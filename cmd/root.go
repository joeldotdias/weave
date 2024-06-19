package cmd

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/config"
	"github.com/joeldotdias/weave/internal/tui/multiInput"
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
	// fmt.Printf("Add all: %t\n", conf.Add_all)
	// fmt.Printf("Title: %s\n", conf.Title)
	// for k, v := range conf.Symbols {
	// 	fmt.Printf("Name: %s, Emoji: %s\n", k, v)
	// }

	opts := Opts{Title: &textInput.Response{}, Symbol: &multiInput.Selected{}}

	tprogram := tea.NewProgram(textInput.InitTextInputModel(opts.Title, "Your title here..."))

	var title string
	if len(conf.Title) != 0 {
		title = conf.Title
	} else {
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(err)
		}
		title = opts.Title.Value()
	}

	fmt.Println("\nYour title: ", title)

	multiOptions := SymbolsList{
		options: conf.SymbolChoices(),
	}

	tprogram = tea.NewProgram(multiInput.InitMultiInputModel(multiOptions.options, opts.Symbol, "Choose symbol"))
	if _, err := tprogram.Run(); err != nil {
		cobra.CheckErr(err)
	}
	symbol := opts.Symbol.Value()
	fmt.Println("\nYour symbol: ", symbol)

}

func Execute() {
	bindFlags()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

type Opts struct {
	Title  *textInput.Response
	Symbol *multiInput.Selected
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
