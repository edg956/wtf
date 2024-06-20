package fake_mail

import (
	"github.com/wtfutil/wtf/modules/mail/model"
)

type Client struct {
	mailboxes []model.Mailbox
}

func (c *Client) GetMailboxes() ([]model.Mailbox, error) {
	return c.mailboxes, nil
}

func (c *Client) GetMessages(mailbox model.Mailbox) ([]model.Envelope, error) {
	return mailbox.Messages, nil
}

var client *Client = nil

func NewClient() *Client {
	if client == nil {
		client = &Client{
			mailboxes: []model.Mailbox{
				{
					Name: "Inbox",
					Messages: []model.Envelope{
						{
							From: model.Contact{
								Name:  "John Doe",
								Email: "john@doe.com",
							},
							To:      []model.Contact{{Name: "Jane Doe", Email: "jane@doe.com"}},
							Subject: "Hello",
							Content: &model.Message{
								Body: "Hi Jane, how are you?\n\nJohn",
							},
						},
					},
				},
				{
					Name: "Sent",
				},
				{
					Name: "Drafts",
				},
			},
		}
	}
	return client
}
