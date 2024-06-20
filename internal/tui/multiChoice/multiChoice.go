package multiChoice

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/tui"
)

type Selected struct {
	Choice string
}

func (s *Selected) Update(value string) {
	s.Choice = strings.TrimSpace(value)
}

func (s *Selected) Value() string {
	return s.Choice
}

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
	choice   *Selected
	header   string
	exit     *bool
}

func InitMultiChoiceModel(choices []string, selected *Selected, header string) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		choice:   selected,
		header:   tui.HeaderStyle.Render(header),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if len(m.selected) == 1 {
				m.selected = make(map[int]struct{})
			}
			_, exists := m.selected[m.cursor]
			if exists {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}

		case "y":
			if len(m.selected) == 1 {
				for selectedKey := range m.selected {
					m.choice.Update(m.choices[selectedKey])
					m.cursor = selectedKey
				}
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "\n" + m.header + "\n"

	var option string
	for i, choice := range m.choices {
		cursor := " "
		option = choice
		if m.cursor == i {
			cursor = tui.FocusedStyle.Render(">")
			option = tui.SelectedItemStyle.Render(choice)
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = tui.FocusedStyle.Render("x")
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, option)
	}

	s += fmt.Sprintf("\nPress %s to confirm choice\n", tui.FocusedStyle.Render("y"))

	return s
}
