package views

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/snwzt/raccoon/internal/models"
)

type model struct {
	spinner       spinner.Model
	cancelContext context.CancelFunc
	status        []*models.Status
	err           error
}

func InitialModel(cancelContext context.CancelFunc, status []*models.Status) model {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner:       spin,
		cancelContext: cancelContext,
		status:        status,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.cancelContext()
			return m, tea.Quit
		} else {
			return m, nil
		}

	case error:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var builder strings.Builder

	for _, data := range m.status {
		curSize := int64(0)
		for _, size := range data.Parts {
			curSize += size
		}

		if data.Done && data.Err != nil {
			builder.WriteString(fmt.Sprintf(" %s [FAILED] \n", data.Path))
		} else if data.Done && data.Err == nil {
			builder.WriteString(fmt.Sprintf(" %s [DONE] \n", data.Path))
		} else {
			builder.WriteString(fmt.Sprintf(" %s %s [%0.1f MB/%0.1f MB] [%0.1f%%] \n", m.spinner.View(), data.Path,
				float64(data.FinalSize)/(1024*1024),
				float64(curSize)/(1024*1024),
				(float64(curSize)/float64(data.FinalSize))*100))
		}
	}

	return builder.String()
}
