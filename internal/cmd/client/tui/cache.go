package tui

import (
	"fmt"
	"slices"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/lordvidex/gostream/internal/entity"
)

type cacheModel struct {
	li list.Model
	cl *Client
}

// Update implements tea.Model.
func (m *cacheModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			item := m.li.SelectedItem().(*cacheItem)
			return m, m.cl.DeletePet(item.Id)
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}
	m.li, cmd = m.li.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m *cacheModel) View() string {
	return m.li.View()
}

func (m *cacheModel) DeleteItem(item uint64) {
	idx := slices.IndexFunc(m.li.Items(), func(i list.Item) bool {
		return i.(*cacheItem).Id == item
	})
	if idx == -1 {
		return
	}
	m.li.RemoveItem(idx)
}

func (m *cacheModel) SetSize(width, height int) {
	m.li.SetSize(width, height)
}

func (m *cacheModel) SetItems(items []entity.Pet) {
	its := make([]list.Item, len(items))
	for i, item := range items {
		its[i] = &cacheItem{item}
	}
	if cmd := m.li.SetItems(its); cmd != nil {
		m.li, _ = m.li.Update(cmd())
	}
}

func (m *cacheModel) AddItem(pet entity.Pet) {
	it := &cacheItem{pet}
	idx := slices.IndexFunc(m.li.Items(), func(i list.Item) bool {
		return i.(*cacheItem).Id == pet.Id
	})
	if idx == -1 {
		m.li.InsertItem(len(m.li.Items()), it)
		return
	}
	if cmd := m.li.SetItem(idx, it); cmd != nil {
		m.li, _ = m.li.Update(cmd())
	}
}

func newCache(cl *Client) *cacheModel {
	li := list.New([]list.Item{}, list.NewDefaultDelegate(), 80, 80)
	li.Title = "local cache"
	return &cacheModel{
		li: li,
		cl: cl,
	}
}

func (m cacheModel) Init() tea.Cmd {
	return nil
}

type cacheItem struct {
	entity.Pet
}

func (i cacheItem) FilterValue() string { return i.Name }
func (i cacheItem) Title() string       { return fmt.Sprintf("pet %d", i.Id) }
func (i cacheItem) Description() string { return fmt.Sprintf("%s: %s, age: %d", i.Kind, i.Name, i.Age) }
