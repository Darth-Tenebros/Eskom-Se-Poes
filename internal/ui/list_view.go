package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"regexp"
)

type Item struct {
	AreaName string
}

func AreasToListItems(areas []string) []list.Item {
	var items []list.Item
	for _, area := range areas {
		item := Item{
			AreaName: area,
		}
		items = append(items, item)
	}
	return items
}

func (i Item) Title() string       { return i.AreaName }
func (i Item) Description() string { return "random desc" }
func (i Item) FilterValue() string { return i.AreaName }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListModel struct {
	List         list.Model
	SelectedItem int
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			areaName := m.List.Items()[m.List.Index()]
			re := regexp.MustCompile("[^a-zA-Z0-9-]")
			return TableModel{Table: LoadTableView(re.ReplaceAllString(areaName.FilterValue(), ""))}, nil
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
func (m ListModel) View() string {
	m.List.SetHeight(35)
	m.List.SetWidth(50)
	return docStyle.Render(m.List.View())
}
