package tui

import tea "github.com/charmbracelet/bubbletea"

type editModel struct {
	// form *huh.Form
}

// Update implements tea.Model.
func (m *editModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	// TODO:
	return m, nil
}

// View implements tea.Model.
func (m *editModel) View() string {
	// TODO:
	return "edit"
}

func newEdit() *editModel {
	return &editModel{}
}

func (m editModel) Init() tea.Cmd {
	return nil
}

