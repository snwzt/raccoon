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
	quitting      *bool
}

func InitialModel(cancelContext context.CancelFunc, status []*models.Status) model {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	var quit bool

	return model{
		spinner:       spin,
		cancelContext: cancelContext,
		status:        status,
		quitting:      &quit,
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
		}
		return m, nil

	case error:
		m.err = msg
		return m, nil

	default:
		if *m.quitting {
			m.cancelContext()
			return m, tea.Quit
		}

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

	quit := true

	for _, data := range m.status {
		curSize := int64(0)
		for _, size := range data.Parts {
			curSize += size
		}

		if data.Done && data.Err != nil {
			builder.WriteString(fmt.Sprintf(" %s [FAILED] \n", data.Name))
		} else if data.Done && data.Err == nil {
			builder.WriteString(fmt.Sprintf(" %s [DONE] \n", data.Name))
		} else {
			quit = false
			builder.WriteString(fmt.Sprintf(" %s%s [%0.1f MB/%0.1f MB] [%0.1f%%] \n", m.spinner.View(), data.Name,
				float64(data.FinalSize)/(1024*1024),
				float64(curSize)/(1024*1024),
				(float64(curSize)/float64(data.FinalSize))*100))
		}
	}

	*m.quitting = quit

	return builder.String()
}
