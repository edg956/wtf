package pages

import (
	"github.com/rivo/tview"
	mailclient "github.com/wtfutil/wtf/modules/mail/client"
	"github.com/wtfutil/wtf/modules/mail/model"
)

type MailboxesPage struct {
	Page
	mailboxes []model.Mailbox
	index     int
	error     error
}

func (m *MailboxesPage) Render(renderFunc func([]string) string) (string, string) {
	if m.error != nil {
		return "", m.error.Error()
	}

	var lines = make([]string, len(m.mailboxes))
	for i, mailbox := range m.mailboxes {
		lines[i] = mailbox.Name
	}

	return "", renderFunc(lines)
}

func (m *MailboxesPage) Select(index int) IPage {
	if index < 0 || index >= len(m.mailboxes) {
		return nil
	}
	return NewMailboxPage(m.view, m.client, m.mailboxes[index])
}

func (m *MailboxesPage) NumberOfItems() int {
	return len(m.mailboxes)
}

func NewMailboxesPage(view *tview.TextView, client mailclient.MailClient) *MailboxesPage {
	mailboxes, err := client.GetMailboxes()

	return &MailboxesPage{
		Page: Page{
			view:   view,
			client: client,
		},
		mailboxes: mailboxes,
		index:     -1,
		error:     err,
	}
}
