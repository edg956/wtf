package mail

import (
	"github.com/emersion/go-imap"
)

type clientSettings struct{}

type IMAPClient struct {
	client goImapClient
}

type goImapClient interface {
	List(ref, name string, mailboxes chan *imap.MailboxInfo) error
	Status(name string, items []imap.StatusItem) (*imap.MailboxStatus, error)
}

func newImapClient(settings *clientSettings) *IMAPClient {
	return &IMAPClient{}
}

func (client *IMAPClient) listMailboxes(config *pageConfig) ([]*mailbox, error) {
	mailboxes := make(chan *imap.MailboxInfo, config.pageSize)
	done := make(chan error, 1)
	defer close(done)

	go func() {
		done <- client.client.List("", "*", mailboxes)
	}()

	if err := <-done; err != nil {
		return nil, err
	}

	mailboxesArray := make([]*mailbox, 0, len(mailboxes))

	for mbox := range mailboxes {
		mailboxStatus, err := client.client.Status(mbox.Name, []imap.StatusItem{imap.StatusMessages, imap.StatusUnseen})
		if err != nil {
			mailboxStatus = &imap.MailboxStatus{Unseen: 0, Messages: 0}
		}
		mailboxesArray = append(mailboxesArray, &mailbox{
			name:     mbox.Name,
			unread:   int(mailboxStatus.Unseen),
			messages: int(mailboxStatus.Messages),
		})
	}

	return mailboxesArray, nil
}
