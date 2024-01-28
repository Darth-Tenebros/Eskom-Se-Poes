package utils

import (
	"eskom-se-poes/internal/ui"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"net/http"
)

func LoadTableView(areaName string) tea.Model {
	link := "https://eskom-calendar-api.shuttleapp.rs/outages/"
	// TODO: FETCH FROM COMMAND LINE
	location := "city-of-cape-town-area-15"

	schedule, err := getSchedule(link, location)
	if err != nil {
		fmt.Print(err)
	}

	var rows []table.Row
	for _, outage := range schedule.Times {
		rows = append(rows, outage.OutageToSlice())
	}

	columns := []table.Column{
		{Title: "Stage", Width: 5},
		{Title: "Area Name", Width: 30},
		{Title: "Start", Width: 30},
		{Title: "Finish", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := ui.TableModel{Table: t}
	return m
}

func getSchedule(link, area string) (*Schedule, error) {
	fullUrl := link + area
	var schedule Schedule

	res, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = schedule.UnmarshalResponse(body)
	if err != nil {
		return nil, err
	}

	return &schedule, nil
}

func AreasToListItems(areas []string) []list.Item {
	var items []list.Item
	for _, area := range areas {
		item := ui.Item{
			AreaName: area,
		}
		items = append(items, item)
	}
	return items
}
