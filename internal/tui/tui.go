package tui

import "github.com/charmbracelet/lipgloss"

var (
	HeaderStyle       = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#1e1e1e")).Bold(true).Padding(0, 1, 0)
	HighlightStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#344CEB")).Bold(true)
	HelpStyle         = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#888a89"))
	ErrStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	NotifStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#A351FE")).Bold(true)
)
