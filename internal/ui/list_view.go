package ui

import (
	"log"
	"regexp"
	"watts-up/internal/repository"
	"watts-up/internal/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (i Item) Title() string {
	re := regexp.MustCompile("[^a-zA-Z0-9-]")
	return re.ReplaceAllString(i.AreaName, "")
}
func (i Item) Description() string { return "" }
func (i Item) FilterValue() string { return i.Title() }

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListModel struct {
	List         list.Model
	SelectedItem int
	Infavourites bool
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
			return TableModel{Table: LoadTableView(areaName.FilterValue())}, nil

		case "a":
			areaName := m.List.Items()[m.List.Index()]
			re := regexp.MustCompile("[^a-zA-Z0-9-]")

			db, err := utils.OpenDatabase()
			if err != nil {
				return m, nil
			}
			defer db.Close()

			areaRepo := repository.NewAreaRepo(db)
			_, err = areaRepo.AddFavourite(re.ReplaceAllString(areaName.FilterValue(), ""))
			if err != nil {
				return m, nil
			}
			return m, nil

		case "b":
			if m.Infavourites {
				m.Infavourites = false
			}
			listModel := LoadListView()
			return listModel, nil

		case "f":
			db, err := utils.OpenDatabase()
			if err != nil {
				return m, nil
			}
			defer db.Close()
			areaRepo := repository.NewAreaRepo(db)

			names, err := areaRepo.GetAllAreaNames()
			if err != nil {
				log.Println(err)
				return m, nil
			}
			items := AreasToListItems(names)
			m := ListModel{List: list.New(items, list.NewDefaultDelegate(), 10, 10)}
			m.List.Title = "SA Load Shedding Areas - FAVOURITES"
			m.Infavourites = true
			return m, nil

		case "d":
			if !m.Infavourites {
				return m, nil
			}

			// delete the area
			db, err := utils.OpenDatabase()
			if err != nil {
				log.Println(err)
				return m, nil
			}
			defer db.Close()
			areaRepo := repository.NewAreaRepo(db)
			id := m.List.Items()[m.List.Index()].FilterValue()
			err = areaRepo.DeleteAreaFromFavourites(utils.GenerateId(id))
			if err != nil {
				log.Println(err)
			}

			// fetch current db data, re-render ui
			// TODO: extract this to a util function
			names, err := areaRepo.GetAllAreaNames()
			if err != nil {
				log.Println(err)
				return m, nil
			}
			items := AreasToListItems(names)
			m := ListModel{List: list.New(items, list.NewDefaultDelegate(), 10, 10)}
			m.List.Title = "SA Load Shedding Areas - FAVOURITES"
			m.Infavourites = true
			return m, nil

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
