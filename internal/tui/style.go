package tui

import "github.com/charmbracelet/lipgloss"

var (
	HeaderStyle       = lipgloss.NewStyle().Background(lipgloss.Color("#01FAC6")).Foreground(lipgloss.Color("#1e1e1e")).Bold(true).Padding(0, 1, 0)
	FocusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#37f7fa"))
	HelpStyle         = lipgloss.NewStyle().Align().Foreground(lipgloss.Color("#888a89"))
)
