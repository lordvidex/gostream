package tui

import (
	"github.com/charmbracelet/bubbles/timer"
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
	BorderStyle(lipgloss.RoundedBorder()).
	Padding(1)
var selected = box.
	BorderForeground(lipgloss.Color("#ff00ff"))

type homeModel struct {
	cl       *Client
	models   []tea.Model
	loaded   bool
	quitting bool
	selected section
	width    int
	height   int
}

func newHome(cl *Client) homeModel {
	return homeModel{
		cl: cl,
		models: []tea.Model{
			newLogs(),
			newCache(cl),
			newEdit(cl),
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
			return selected.Render(lipgloss.Place(width, height-5, lipgloss.Left, lipgloss.Top, m.models[s].View()))
		}
		return box.Render(lipgloss.Place(width-5, height-5, lipgloss.Left, lipgloss.Top, m.models[s].View()))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			styled(cache, m.width*3/4, (m.height/2)-5),
			styled(logs, m.width*3/4, (m.height/2)-5),
		),
		styled(edit, m.width/4, m.height),
	)
}

func (m homeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - box.GetHorizontalFrameSize()
		m.height = msg.Height - box.GetVerticalFrameSize()
		m.models[logs], cmd = m.models[logs].Update(tea.WindowSizeMsg{Width: m.width / 2, Height: m.height / 2})
		cmds = append(cmds, cmd)
		m.models[cache], cmd = m.models[cache].Update(tea.WindowSizeMsg{Width: m.width / 2, Height: m.height / 2})
		cmds = append(cmds, cmd)
		m.models[edit], cmd = m.models[edit].Update(tea.WindowSizeMsg{Width: m.width / 2, Height: m.height})
		cmds = append(cmds, cmd)
		m.loaded = true
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.selected = (m.selected + 1) % section(len(m.models))
			return m, m.models[edit].Init()
		}
	case restartEditMsg:
		m.models[edit] = newEdit(m.cl)
		return m, m.models[edit].Init()
	case timer.TimeoutMsg, timer.TickMsg: // always propagate to edit timer
		m.models[edit], cmd = m.models[edit].Update(msg)
		cmds = append(cmds, cmd)
	case updateMsg:
		m.models[cache].(*cacheModel).AddItem(msg.Pet)
	case deleteMsg:
		m.models[cache].(*cacheModel).DeleteItem(msg.ID)
	case snapshotMsg:
		m.models[cache].(*cacheModel).SetItems(msg.Pets)
	case logMsg:
		m.models[logs].(*logsModel).AddItem(msg.Source, msg.JSON)
		return m, nil
	}

	m.models[m.selected], cmd = m.models[m.selected].Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// all listened logs
type logMsg struct {
	JSON   string
	Source string
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
