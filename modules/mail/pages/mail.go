package pages

import (
	"fmt"
	"github.com/rivo/tview"
	mailclient "github.com/wtfutil/wtf/modules/mail/client"
	"github.com/wtfutil/wtf/modules/mail/model"
)

type MailPage struct {
	Page
	message model.Envelope
	error   error
}

func (m *MailPage) Render(renderFunc func([]string) string) (string, string) {
	if m.error != nil {
		return "", m.error.Error()
	}

	return "", fmt.Sprintf(
		"From: %s\nSubject: %s\n\n%s",
		m.message.From.String(),
		m.message.Subject,
		m.message.Content.Body,
	)
}

func (m *MailPage) NumberOfItems() int {
	return 0
}

func (m *MailPage) Select(index int) IPage {
	return nil
}

func NewMailPage(view *tview.TextView, client mailclient.MailClient, message model.Envelope) *MailPage {
	return &MailPage{
		Page: Page{
			view:   view,
			client: client,
		},
		message: message,
	}
}
