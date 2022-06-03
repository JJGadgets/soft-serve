package repo

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/soft-serve/git"
	"github.com/charmbracelet/soft-serve/ui/styles"
	"github.com/muesli/reflow/truncate"
)

// LogItem is a item in the log list that displays a git commit.
type LogItem struct {
	*git.Commit
}

// ID implements selector.IdentifiableItem.
func (i LogItem) ID() string {
	return i.Commit.ID.String()
}

// Title returns the item title. Implements list.DefaultItem.
func (i LogItem) Title() string {
	if i.Commit != nil {
		return strings.Split(i.Commit.Message, "\n")[0]
	}
	return ""
}

// Description returns the item description. Implements list.DefaultItem.
func (i LogItem) Description() string { return "" }

// FilterValue implements list.Item.
func (i LogItem) FilterValue() string { return i.Title() }

// LogItemDelegate is the delegate for LogItem.
type LogItemDelegate struct {
	style *styles.Styles
}

// Height returns the item height. Implements list.ItemDelegate.
func (d LogItemDelegate) Height() int { return 1 }

// Spacing returns the item spacing. Implements list.ItemDelegate.
func (d LogItemDelegate) Spacing() int { return 0 }

// Update updates the item. Implements list.ItemDelegate.
func (d LogItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

// Render renders the item. Implements list.ItemDelegate.
func (d LogItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(LogItem)
	if !ok {
		return
	}
	if i.Commit == nil {
		return
	}

	hash := i.Commit.ID.String()
	leftMargin := d.style.LogItemSelector.GetMarginLeft() +
		d.style.LogItemSelector.GetWidth() +
		d.style.LogItemHash.GetMarginLeft() +
		d.style.LogItemHash.GetWidth() +
		d.style.LogItemInactive.GetMarginLeft()
	title := truncateString(i.Title(), m.Width()-leftMargin, "…")
	if index == m.Index() {
		fmt.Fprint(w, d.style.LogItemSelector.Render(">")+
			d.style.LogItemHash.Bold(true).Render(hash[:7])+
			d.style.LogItemActive.Render(title))
	} else {
		fmt.Fprint(w, d.style.LogItemSelector.Render(" ")+
			d.style.LogItemHash.Render(hash[:7])+
			d.style.LogItemInactive.Render(title))
	}
}

func truncateString(s string, max int, tail string) string {
	if max < 0 {
		max = 0
	}
	return truncate.StringWithTail(s, uint(max), tail)
}
