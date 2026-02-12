package picker

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// ErrCancelled is returned when the user cancels the picker (Ctrl+C / Esc).
var ErrCancelled = errors.New("picker cancelled")

// entry wraps a key/path pair as a list.Item.
type entry struct {
	key  string
	path string
}

func (e entry) FilterValue() string { return e.key }
func (e entry) Title() string       { return e.key }
func (e entry) Description() string { return e.path }

type model struct {
	list     list.Model
	selected *entry
	quitting bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if e, ok := m.list.SelectedItem().(entry); ok {
				m.selected = &e
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	help := "\033[2m↑/↓ navigate  /  filter  enter select  esc cancel\033[0m"
	return fmt.Sprintf("\n%s\n\n%s", m.list.View(), help)
}

// Pick presents an interactive filterable list of key/path pairs.
// Returns the path of the selected entry, or ErrCancelled if the user exits.
func Pick(keys []string, pathByKey map[string]string) (string, error) {
	items := make([]list.Item, len(keys))
	for i, k := range keys {
		items[i] = entry{key: k, path: pathByKey[k]}
	}

	height := len(keys) + 6
	if height > 20 {
		height = 20
	}

	l := list.New(items, list.NewDefaultDelegate(), 70, height)
	l.Title = "仓库选择"
	l.SetFilteringEnabled(true)
	l.SetShowStatusBar(false)

	m := model{list: l}
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return "", err
	}

	final := result.(model)
	if final.quitting || final.selected == nil {
		return "", ErrCancelled
	}
	return final.selected.path, nil
}
