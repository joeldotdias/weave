package multiInput

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selected struct {
	Choice string
}

func (s *Selected) Update(value string) {
	s.Choice = value
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

func InitMultiInputModel(choices []string, selected *Selected, header string) model {
	var header_style = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#030303")).Bold(true).Padding(0, 1, 0)
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		choice:   selected,
		header:   header_style.Render(header),
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
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	s := m.header + "\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == 1 {
			cursor = focusedStyle.Render(">")
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = focusedStyle.Render("x")
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += fmt.Sprintf("\nPress %s to confirm choice\n", focusedStyle.Render("y"))

	return s
}
