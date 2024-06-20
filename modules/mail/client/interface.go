package client

import (
	"github.com/wtfutil/wtf/modules/mail/model"
)

type MailClient interface {
	GetMailboxes() ([]model.Mailbox, error)
	GetMessages(mailbox model.Mailbox) ([]model.Envelope, error)
}
