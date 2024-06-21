package textArea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/config"
	"github.com/joeldotdias/weave/internal/tui"
)

type errMsg error

type Description struct {
	content string
}

func (d *Description) Value() string {
	return d.content
}

func (d *Description) update(input string) {
	d.content = input
}

type model struct {
	textArea    textarea.Model
	err         error
	description *Description
	header      string
	hint        string
	colors      *config.Colors
	exit        *bool
}

func InitTextAreaModel(description *Description, header string, hint string, colors *config.Colors, exit *bool) model {
	ta := textarea.New()
	ta.Focus()

	return model{
		textArea:    ta,
		err:         nil,
		description: description,
		header:      header,
		hint:        hint,
		colors:      colors,
		exit:        exit,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			*m.exit = true
			return m, tea.Quit
		case tea.KeyCtrlY:
			input := m.textArea.Value()
			if len(input) > 0 && len(strings.Fields(input)) <= 72 {
				m.description.update(input)
				m.textArea.Blur()
				return m, tea.Quit
			}
		default:
			if !m.textArea.Focused() {
				cmd = m.textArea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case errMsg:
		m.err = msg
		*m.exit = true
		return m, nil
	}

	m.textArea, cmd = m.textArea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n\n%s\n%s",
		tui.HeaderStyle(m.colors.HeaderBg, m.colors.HeaderFg).Render(m.header),
		m.textArea.View(),
		tui.HelpStyle.Render("("+m.hint+")"),
		tui.HelpStyle.Render(fmt.Sprintf("Press %s to confirm and make the commit", tui.HighlightStyle.Render("ctrl+y"))),
	)
}
