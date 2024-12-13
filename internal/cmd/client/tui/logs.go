package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type logsModel struct {
	li list.Model
}

func newLogs() *logsModel {
	li := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	li.SetShowTitle(false)
	return &logsModel{
		li: li,
	}
}

func (m *logsModel) Init() tea.Cmd {
	return nil
}

func (m *logsModel) SetSize(width, height int) {
	m.li.SetSize(width, height)
}

func (m *logsModel) View() string {
	return m.li.View()
}

func (m *logsModel) AddItem(entry string) {
	m.li.InsertItem(len(m.li.Items()), &serverLog{entry})
}

func (m *logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.li, cmd = m.li.Update(msg)
	return m, cmd
}

type serverLog struct {
	Entry string
}

func (l *serverLog) FilterValue() string { return l.Entry }
func (l *serverLog) Title() string       { return "" }
func (l *serverLog) Description() string {
	return "some very logn\n stasdkfadsfasdgardgdgfdshgjreklgherjgohewrqgk;ewqlhgeqrjwghjergerqherqfgjewqgfjeqwohfewjrqgheqrwutfherwjgfheerugbore uvfrevregreg\nasdfjasdfkhdsjgfbdajfkdshafjhhadsfdsagfdsagdgrg"
}
