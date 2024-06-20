package pages

import (
	"github.com/rivo/tview"
	mailclient "github.com/wtfutil/wtf/modules/mail/client"
)

type IPage interface {
	Render(renderFunc func([]string) string) (string, string)
	Select(index int) IPage
	NumberOfItems() int
}

type Page struct {
	view   *tview.TextView
	client mailclient.MailClient
}

func StartPage(view *tview.TextView, client mailclient.MailClient) IPage {
	return NewMailboxesPage(view, client)
}
