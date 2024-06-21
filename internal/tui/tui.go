package tui

import "github.com/charmbracelet/lipgloss"

var (
	HighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	HelpStyle      = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#888a89"))
	ErrStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	NotifStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#A351FE")).Bold(true)
)

const (
	DefaultHeaderBg    string = "#8AF3FFEE"
	DefaultHeaderFg    string = "#1e1e1e"
	DefaultSelectedOpt string = "#3F56EB"
)

func HeaderStyle(headerBg string, headerFg string) lipgloss.Style {
	if headerBg == "" {
		headerBg = DefaultHeaderBg
	}
	if headerFg == "" {
		headerFg = DefaultHeaderFg
	}

	return lipgloss.NewStyle().
		Background(lipgloss.Color(headerBg)).
		Foreground(lipgloss.Color(headerFg)).
		Bold(true).
		Padding(0, 1, 0)
}

func SelectedItemStyle(fg string) lipgloss.Style {
	if fg == "" {
		fg = DefaultSelectedOpt
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Bold(true)
}
