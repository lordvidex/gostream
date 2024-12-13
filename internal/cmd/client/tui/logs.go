package tui

import (
	"bytes"
	"encoding/json"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type logsModel struct {
	li list.Model
	v  viewport.Model
	s  bytes.Buffer
}

func newLogs() *logsModel {
	del := list.NewDefaultDelegate()
	del.ShowDescription = false
	li := list.New([]list.Item{}, del, 0, 0)
	return &logsModel{
		li: li,
		v:  viewport.New(80, 80),
	}
}

func (m *logsModel) Init() tea.Cmd {
	return nil
}

func (m *logsModel) SetSize(width, height int) {
	m.li.SetSize(width, height)
	m.v.Height = height
	m.v.Width = width
}

func (m *logsModel) View() string {
	return m.v.View()
	// return m.li.View()
}

func (m *logsModel) AddItem(entry string) {
	json.Indent(&m.s, []byte(entry), "", "	")
	m.s.WriteString("\n")
	m.v.SetContent(m.s.String())
	// m.li.InsertItem(len(m.li.Items()), &serverLog{entry})
}

func (m *logsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.v, cmd = m.v.Update(msg)
	return m, cmd
}

type serverLog struct {
	Entry string
}

func (l *serverLog) FilterValue() string { return "" }
func (l *serverLog) Title() string       { return l.Entry }
func (l *serverLog) Description() string {
	return ""
}
