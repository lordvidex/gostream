package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lordvidex/gostream/internal/entity"
)

type section int

const (
	logs section = iota
	cache
	edit
)

var box = lipgloss.NewStyle().
	Padding(0, 1)
var selected = box.
	BorderStyle(lipgloss.BlockBorder()).
	BorderForeground(lipgloss.Color("#ff00ff"))

type homeModel struct {
	models   []tea.Model
	loaded   bool
	quitting bool
	selected section
	width    int
	height   int
}

func newHome() homeModel {
	return homeModel{
		models: []tea.Model{
			newLogs(),
			newCache(),
			newEdit(),
		},
	}
}

func (m homeModel) Init() tea.Cmd {
	return nil
}

func (m homeModel) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading ..."
	}

	styled := func(s section, width, height int) string {
		if s == m.selected {
			return selected.Render(lipgloss.Place(width, height, lipgloss.Left, lipgloss.Top, m.models[s].View()))
		}
		return box.Render(lipgloss.Place(width, height, lipgloss.Left, lipgloss.Top, m.models[s].View()))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		lipgloss.JoinVertical(
			lipgloss.Left,
			styled(cache, m.width/2, m.height/2),
			styled(logs, m.width/2, m.height/2),
		),
		styled(edit, m.width/2, m.height),
	)
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width -
			(box.GetHorizontalFrameSize() * 4)
		m.height = msg.Height - (box.GetVerticalFrameSize() * 4)
		m.models[logs].(*logsModel).SetSize(m.width/2, m.height/2)
		m.loaded = true
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.selected = (m.selected + 1) % section(len(m.models))
			return m, nil
		}
	// case updateMsg:
	// case deleteMsg:
	// case snapshotMsg:
	case logMsg:
		m.models[logs].(*logsModel).AddItem(msg.JSON)
		return m, nil
	}
	var cmd tea.Cmd
	m.models[m.selected], cmd = m.models[m.selected].Update(msg)
	return m, cmd
}

// all listened logs
type logMsg struct {
	JSON string
}

// only pet updates
type updateMsg struct {
	Pet entity.Pet
}

// only pet snapshots
type snapshotMsg struct {
	Pets []entity.Pet
}

// only pet deletes
type deleteMsg struct {
	ID uint64
}
