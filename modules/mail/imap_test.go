package mail

import (
	"fmt"
	"github.com/emersion/go-imap"
	"testing"
)

type fakeGoImapClient struct {
	mailboxInfos    []*imap.MailboxInfo
	mailboxStatuses map[string]*imap.MailboxStatus
	error           error
}

func (client *fakeGoImapClient) List(ref, name string, mailboxes chan *imap.MailboxInfo) error {
	defer close(mailboxes)

	if client.error != nil {
		return client.error
	}

	for _, mailboxInfo := range client.mailboxInfos {
		mailboxes <- mailboxInfo
	}

	return nil
}

func (client *fakeGoImapClient) Status(name string, items []imap.StatusItem) (*imap.MailboxStatus, error) {
	return client.mailboxStatuses[name], nil
}

func TestIMAPListMailboxes(t *testing.T) {
	t.Run("returns a list of mailboxes", func(t *testing.T) {
		mailboxInfos := []*imap.MailboxInfo{
			&imap.MailboxInfo{
				Name: "INBOX",
			},
			&imap.MailboxInfo{
				Name: "Newsletter",
			},
		}
		mailboxStatuses := map[string]*imap.MailboxStatus{
			"INBOX": &imap.MailboxStatus{
				Name:     "INBOX",
				Unseen:   1,
				Messages: 10,
			},
			"Newsletter": &imap.MailboxStatus{
				Name:     "Newsletter",
				Unseen:   1,
				Messages: 10,
			},
		}

		client := &IMAPClient{
			client: &fakeGoImapClient{
				mailboxInfos:    mailboxInfos,
				mailboxStatuses: mailboxStatuses,
			},
		}

		result, err := client.listMailboxes(&pageConfig{
			page:     0,
			pageSize: 10,
		})

		if err != nil {
			t.Errorf("Expected no error, got %s", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 mailboxes, got %d", len(result))
		}

		expectedInbox := mailbox{
			name:     "INBOX",
			unread:   1,
			messages: 10,
		}
		expectedNewsletter := mailbox{
			name:     "Newsletter",
			unread:   1,
			messages: 10,
		}

		if *result[0] != expectedInbox {
			t.Errorf("Expected mailbox to be %v, got %v", expectedInbox, result[0])
		}

		if *result[1] != expectedNewsletter {
			t.Errorf("Expected mailbox to be %v, got %v", expectedNewsletter, result[1])
		}
	})
	t.Run("returns an error message when the client returns an error", func(t *testing.T) {
		client := &IMAPClient{
			client: &fakeGoImapClient{
				error: fmt.Errorf("error"),
			},
		}

		_, err := client.listMailboxes(&pageConfig{
			page:     0,
			pageSize: 10,
		})

		if err == nil {
			t.Error("Expected an to be returned error")
		}
	})
}
