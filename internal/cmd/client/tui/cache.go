package tui

import tea "github.com/charmbracelet/bubbletea"

type cacheModel struct{}

// Update implements tea.Model.
func (m *cacheModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	// TODO:
	return m, nil
}

// View implements tea.Model.
func (m *cacheModel) View() string {
	return "cache"
}

func newCache() *cacheModel {
	return &cacheModel{}
}

func (m cacheModel) Init() tea.Cmd {
	return nil
}
