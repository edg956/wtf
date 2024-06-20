package pages

import (
	"fmt"
	"github.com/rivo/tview"
	mailclient "github.com/wtfutil/wtf/modules/mail/client"
	"github.com/wtfutil/wtf/modules/mail/model"
)

type MailboxPage struct {
	Page
	mailbox  model.Mailbox
	messages []model.Envelope
	index    int
	error    error
}

func (m *MailboxPage) Render(renderFunc func([]string) string) (string, string) {
	if m.error != nil {
		return "", m.error.Error()
	}

	var lines = make([]string, len(m.messages))
	for i, mailbox := range m.messages {
		lines[i] = fmt.Sprintf("%s: %s", mailbox.From.String(), mailbox.Subject)
	}
	return m.mailbox.Name, renderFunc(lines)
}

func (m *MailboxPage) NumberOfItems() int {
	return len(m.messages)
}

func (m *MailboxPage) Select(index int) IPage {
	if index < 0 || index >= len(m.messages) {
		return nil
	}
	return NewMailPage(m.view, m.client, m.messages[index])
}

func NewMailboxPage(view *tview.TextView, client mailclient.MailClient, mailbox model.Mailbox) *MailboxPage {
	messages, err := client.GetMessages(mailbox)

	return &MailboxPage{
		Page: Page{
			view:   view,
			client: client,
		},
		mailbox:  mailbox,
		messages: messages,
		index:    -1,
		error:    err,
	}
}
