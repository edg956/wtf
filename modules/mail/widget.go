package mail

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

type imapClient interface {
	listMailboxes(config *pageConfig) ([]*mailbox, error)
}

type mailbox struct {
	name     string
	unread   int
	messages int
}

type pageConfig struct {
	page     int
	pageSize int
}

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget

	settings   *Settings
	imapClient imapClient
	pageConfig *pageConfig
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.common),

		settings:   settings,
		imapClient: newImapClient(nil),
	}

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {

	// The last call should always be to the display function
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) listMailboxes() string {
	mailboxes, err := widget.imapClient.listMailboxes(widget.pageConfig)

	if err != nil {
		return "Error loading mailboxes"
	}

	content := ""
	for _, mbox := range mailboxes {
		content += fmt.Sprintf("%s (unread: %d/messages: %d)\n", mbox.name, mbox.unread, mbox.messages)
	}

	return content
}

func (widget *Widget) content() string {
	return widget.listMailboxes()
}

func (widget *Widget) display() {
	widget.Redraw(func() (string, string, bool) {
		return widget.CommonSettings().Title, widget.content(), false
	})
}
