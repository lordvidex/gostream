package tui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff00ff"))

type editModel struct {
	form     *huh.Form
	cl       *Client
	t        timer.Model
	finished bool
}

func (m *editModel) createPet() tea.Cmd {
	age, _ := strconv.Atoi(m.form.GetString("age"))

	var cmds []tea.Cmd
	cmds = append(cmds, m.cl.CreatePet(
		m.form.GetString("name"),
		m.form.GetString("kind"),
		uint32(age),
	)) // make the request

	cmds = append(cmds, m.t.Init()) // start the timer

	return tea.Batch(cmds...)
}

// Update implements tea.Model.
func (m *editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.finished && m.form.State == huh.StateCompleted {
		m.t = timer.New(time.Second * 5)
		m.finished = true
		return m, m.createPet()
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case timer.TimeoutMsg:
		return m, m.restartEdit
	case timer.TickMsg:
		m.t, cmd = m.t.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View implements tea.Model.
func (m *editModel) View() string {
	strs := []string{
		titleStyle.Render("New Pet Form"),
		"\n",
		m.form.View(),
	}
	if m.form.State == huh.StateCompleted {
		strs = append(strs, "You have created a new pet!\n\n")
	}
	if m.finished {
		strs = append(strs, "Reloading form ...\n")
		strs = append(strs, m.t.View())
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		strs...,
	)
}

func (m *editModel) SetSize(width, height int) {
	m.form = m.form.WithWidth(width).WithHeight(height)
}

func (m *editModel) restartEdit() tea.Msg {
	return restartEditMsg{}
}

func newEdit(cl *Client) *editModel {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Pet Name").
				Prompt("?").
				Key("name"),
			huh.NewInput().
				Title("Pet Kind").
				Prompt("?").
				Key("kind"),
			huh.NewInput().
				Title("Pet Age").
				Prompt("?").
				Key("age").
				Validate(func(v string) error {
					_, err := strconv.ParseInt(v, 10, 64)
					return err
				}),
		),
	)
	return &editModel{
		form: form,
		cl:   cl,
	}
}

func (m editModel) Init() tea.Cmd {
	return m.form.Init()
}

type restartEditMsg struct{}
