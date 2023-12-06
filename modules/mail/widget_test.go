package mail

import (
	"fmt"
	"testing"
)

type fakeClient struct {
	mailboxes []*mailbox
	error     error
}

func (client *fakeClient) listMailboxes(config *pageConfig) ([]*mailbox, error) {
	if client.error != nil {
		return nil, client.error
	}
	return client.mailboxes, nil
}

func TestListMailboxes(t *testing.T) {
	t.Run("returns a list of mailboxes", func(t *testing.T) {
		mailboxes := []*mailbox{
			&mailbox{
				name:     "INBOX",
				unread:   1,
				messages: 10,
			},
			&mailbox{
				name:     "Newsletter",
				unread:   1,
				messages: 1,
			},
		}

		widget := &Widget{
			imapClient: &fakeClient{
				mailboxes: mailboxes,
			},
		}

		result := widget.listMailboxes()

		expected := "INBOX (unread: 1/messages: 10)\n" +
			"Newsletter (unread: 1/messages: 1)\n"

		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
	t.Run("returns an error message when the client returns an error", func(t *testing.T) {
		widget := &Widget{
			imapClient: &fakeClient{
				error: fmt.Errorf("error"),
			},
		}

		result := widget.listMailboxes()

		expected := "Error loading mailboxes"

		if result != expected {
			t.Errorf("Expected %s, got %s", expected, result)
		}
	})
}
