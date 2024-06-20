package textInput

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joeldotdias/weave/internal/tui"
)

type errMsg error

type Response struct {
	content string
}

func (r *Response) Value() string {
	return r.content
}

func (r *Response) update(input string) {
	r.content = input
}

type model struct {
	textInput textinput.Model
	err       error
	response  *Response
	header    string
	exit      *bool
}

func InitTextInputModel(response *Response, header string) model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 150
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
		response:  response,
		header:    tui.HeaderStyle.Render(header),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			input := m.textInput.Value()
			if len(input) > 0 {
				m.response.update(input)
				return m, tea.Quit
			}
		case tea.KeyCtrlC, tea.KeyEsc:
			*m.exit = true
			return m, tea.Quit
		}
	case errMsg:
		m.err = msg
		*m.exit = true
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s\n", m.header, m.textInput.View())
}
